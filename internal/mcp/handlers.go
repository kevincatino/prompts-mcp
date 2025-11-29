package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"go.uber.org/zap"

	"prompts-mcp/internal/prompts"
)

type Handlers struct {
	promptsRepo prompts.Repository
	logger      *zap.Logger
}

func NewHandlers(promptsRepo prompts.Repository, logger *zap.Logger) *Handlers {
	return &Handlers{promptsRepo: promptsRepo, logger: logger}
}

type promptInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type listPromptsResult struct {
	Prompts []promptInfo `json:"prompts"`
}

type expandPromptArgs struct {
	Command string `json:"command"`
	Input   string `json:"input"`
}

type expandPromptResult struct {
	Prompt string `json:"prompt"`
}

func (h *Handlers) ListPrompts(ctx context.Context) (ToolResult, error) {
	promptsList, err := h.promptsRepo.ListPrompts(ctx)
	if err != nil {
		return ToolResult{}, err
	}

	result := listPromptsResult{Prompts: make([]promptInfo, 0, len(promptsList))}
	for _, p := range promptsList {
		result.Prompts = append(result.Prompts, promptInfo{
			Name:        p.Name,
			Description: p.Description,
		})
	}

	return buildToolResult(result)
}

func (h *Handlers) ExpandPrompt(ctx context.Context, args expandPromptArgs) (ToolResult, error) {
	if args.Command == "" {
		return ToolResult{}, fmt.Errorf("command is required")
	}

	prompt, err := h.promptsRepo.GetPrompt(ctx, args.Command)
	if err != nil {
		return ToolResult{}, err
	}

	expanded := prompt.Content
	if strings.Contains(expanded, "{{input}}") {
		expanded = strings.ReplaceAll(expanded, "{{input}}", args.Input)
	} else {
		// If no placeholder is present, append the input after two newlines to keep the template intact.
		expanded = expanded + "\n\n" + args.Input
	}

	return buildToolResult(expandPromptResult{Prompt: expanded})
}

func decodeArgs[T any](raw json.RawMessage) (T, error) {
	var args T
	if len(raw) == 0 {
		return args, fmt.Errorf("arguments are required")
	}
	if err := json.Unmarshal(raw, &args); err != nil {
		return args, err
	}
	return args, nil
}

func buildToolResult(structured any) (ToolResult, error) {
	raw, err := json.Marshal(structured)
	if err != nil {
		return ToolResult{}, fmt.Errorf("marshal result: %w", err)
	}

	return ToolResult{
		Content: []Content{
			{Type: "text", Text: string(raw)},
		},
		StructuredContent: structured,
	}, nil
}
