Architecture
============

Overview
--------
- CLI entrypoint (`cmd/prompts/main.go`) parses `--prompts-dir`, builds a zap production logger, validates the directory, and constructs the MCP server.
- The server runs over stdio using JSON-RPC 2.0. It decodes requests, routes to handlers, and encodes responses line-by-line.
- Handlers expose two MCP tools backed by a YAML prompt repository.

Components and Flow
-------------------
1) `cmd/prompts/main.go`
   - Parses `--prompts-dir` flag (must be absolute).
   - Initializes logger via `logging.New()` (`zap.NewProduction()` JSON output).
   - Validates the directory with `validate.Dir`, rejecting relative/root/non-existent/non-dir paths and symlink escapes.
   - Builds `prompts.NewYAMLRepository(promptsDir)` and `mcp.NewServer(logger, repo)`.
   - Starts serving with cancellation on SIGINT/SIGTERM.

2) `internal/validate`
   - `Dir` normalizes and resolves symlinks, ensuring an existing directory that is not `/`.

3) `internal/prompts`
   - `YAMLRepository` lists `.yaml` files in the base directory, trims names for prompt IDs, and loads/validates content.
   - `Prompt.Validate` requires `name` and `content`; descriptions are optional.
   - `loadPrompt` trims whitespace on description/content to keep responses clean.

4) `internal/mcp`
   - `Server.Serve` loops over decoded JSON-RPC requests; returns errors for bad JSON-RPC version or unknown methods.
   - Supported methods: `initialize`, `notifications/initialized` (no-op), `tools/list`, `tools/call`.
   - `listTools` returns the tool metadata for `list_prompts` and `expand_prompt`.
   - `callTool` decodes tool-specific arguments and dispatches to handlers; wraps errors as JSON-RPC errors or MCP tool errors (`ToolResult{isError:true}`).

5) `internal/mcp/handlers`
   - `ListPrompts` returns names/descriptions from the repository.
   - `ExpandPrompt` fetches a prompt by name and either replaces all `{{input}}` occurrences or appends the input separated by two newlines when no placeholder exists.

Data Contracts
--------------
- Input: YAML files shaped as:
  - `description: "short description"`
  - `prompt: |` followed by the prompt body (may include `{{input}}`).
- Tools:
  - `list_prompts` → `{"prompts":[{"name","description"}]}`
  - `expand_prompt` → `{"prompt":"expanded string"}`
- JSON-RPC version enforced: `"jsonrpc":"2.0"`.
- MCP protocol version advertised in `initialize`: `2024-11-05`.

Operational Notes
-----------------
- Logging uses zap production JSON to stdout/stderr; defer `logger.Sync()` in main.
- Server stops cleanly on context cancellation (SIGINT/SIGTERM).
- Prompts directory reading and YAML parsing propagate errors through tool responses.
