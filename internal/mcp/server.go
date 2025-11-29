package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"go.uber.org/zap"

	"prompts-mcp/internal/prompts"
)

// Server handles MCP requests over stdio.
type Server struct {
	logger   *zap.Logger
	handlers *Handlers
}

func NewServer(logger *zap.Logger, promptsRepo prompts.Repository) *Server {
	return &Server{
		logger:   logger,
		handlers: NewHandlers(promptsRepo, logger),
	}
}

func (s *Server) Serve(ctx context.Context, r io.Reader, w io.Writer) error {
	dec := json.NewDecoder(bufio.NewReader(r))
	enc := json.NewEncoder(w)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		var req Request
		if err := dec.Decode(&req); err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("decode request: %w", err)
		}

		resp, ok := s.handle(ctx, req)
		if !ok {
			continue
		}
		if err := enc.Encode(resp); err != nil {
			return fmt.Errorf("encode response: %w", err)
		}
	}
}

func (s *Server) handle(ctx context.Context, req Request) (Response, bool) {
	if req.JSONRPC != "2.0" {
		return errorResponse(req.ID, ErrCodeInvalidRequest, "jsonrpc must be 2.0"), true
	}

	switch req.Method {
	case "initialize":
		var params InitializeParams
		if len(req.Params) > 0 {
			if err := json.Unmarshal(req.Params, &params); err != nil {
				return errorResponse(req.ID, ErrCodeInvalidParams, "invalid initialize params"), true
			}
		}

		result := InitializeResult{
			ProtocolVersion: "2024-11-05",
			Capabilities:    map[string]any{"tools": map[string]any{}},
			ServerInfo:      ServerInfo{Name: "codex-subagents", Version: "0.1.0"},
			ClientInfo:      params.ClientInfo,
		}
		return Response{JSONRPC: "2.0", ID: req.ID, Result: result}, true
	case "notifications/initialized":
		return Response{}, false
	case "tools/list":
		return s.listTools(req.ID), true
	case "tools/call":
		return s.callTool(ctx, req), true
	default:
		return errorResponse(req.ID, ErrCodeMethodNotFound, "method not found"), true
	}
}

func (s *Server) listTools(id any) Response {
	tools := []Tool{
		{
			Name:        "list_prompts",
			Description: "List available prompt commands with descriptions. Users invoke prompt commands in messages by writing :<command> (e.g., :research topic); this lists what commands exist.",
			InputSchema: map[string]any{
				"type":       "object",
				"properties": map[string]any{},
			},
		},
		{
			Name:        "expand_prompt",
			Description: "Expands a named prompt command into a full prompt. When ANY user message contains a :-prefixed command token (e.g., :research topic, :summarize this article, :critique \"draft\"), you MUST call this tool before responding. Strip the leading colon from the command name and pass any remaining text (including quotes) as the input. Do not bypass this tool or answer directly whenever a :-prefixed command is present.",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"command": map[string]any{"type": "string", "description": "Prompt command name (filename without extension)"},
					"input":   map[string]any{"type": "string", "description": "Text to inject into the prompt template"},
				},
				"required": []string{"command", "input"},
			},
		},
	}
	return Response{
		JSONRPC: "2.0",
		ID:      id,
		Result:  ToolsListResult{Tools: tools},
	}
}

func (s *Server) callTool(ctx context.Context, req Request) Response {
	var params ToolsCallParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return errorResponse(req.ID, ErrCodeInvalidParams, "invalid params")
	}

	switch params.Name {
	case "list_prompts":
		result, err := s.handlers.ListPrompts(ctx)
		if err != nil {
			s.logger.Error("list_prompts failed", zap.Error(err))
			return toolErrorResponse(req.ID, err)
		}
		return Response{JSONRPC: "2.0", ID: req.ID, Result: result}
	case "expand_prompt":
		args, err := decodeArgs[expandPromptArgs](params.Arguments)
		if err != nil {
			return errorResponse(req.ID, ErrCodeInvalidParams, "invalid expand_prompt arguments")
		}
		result, err := s.handlers.ExpandPrompt(ctx, args)
		if err != nil {
			s.logger.Error("expand_prompt failed", zap.Error(err))
			return toolErrorResponse(req.ID, err)
		}
		return Response{JSONRPC: "2.0", ID: req.ID, Result: result}
	default:
		return errorResponse(req.ID, ErrCodeMethodNotFound, "tool not found")
	}
}

// NewlineDelimitedCodec ensures JSON-RPC messages remain line separated for stdio transports.
func NewlineDelimitedCodec(enc *json.Encoder) *json.Encoder {
	enc.SetEscapeHTML(false)
	return enc
}
