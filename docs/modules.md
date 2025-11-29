Modules
=======

cmd/prompts
-----------
- `main.go`: CLI entrypoint. Parses `--prompts-dir`, initializes zap logger, validates the directory, constructs YAML repository, builds MCP server, and serves with signal-aware context.

internal/logging
----------------
- `logger.go`: `New()` returns a production zap logger configured for JSON output.

internal/validate
-----------------
- `paths.go`: `Dir(string) (string, error)` ensures a non-root absolute directory exists, resolves symlinks, and rejects files.
- `paths_test.go`: unit coverage for validation cases (empty, relative, root, missing, file path, valid dir).

internal/prompts
----------------
- `model.go`: `Prompt` struct (`Name`, `Description`, `Content`); `Validate` requires `Name` and `Content`.
- `repository.go`: Repository interface defining `ListPrompts` and `GetPrompt`.
- `yaml_repository.go`: `YAMLRepository` implementation that reads `.yaml` files, trims whitespace, validates prompts, and supports listing or fetching by name.

internal/mcp
------------
- `types.go`: JSON-RPC request/response types, MCP tool structures, initialize params/result.
- `errors.go`: helpers for JSON-RPC error responses and MCP tool error wrapping.
- `handlers.go`: tool implementations (`list_prompts`, `expand_prompt`), argument decoding, and utility for building tool results.
- `server.go`: stdio server loop, method routing (`initialize`, `notifications/initialized`, `tools/list`, `tools/call`), tool registry, and logging integration.

examples/prompts
----------------
- `research.yaml`: prompt with `{{input}}` placeholder.
- `summarize.yaml`: prompt with `{{input}}` placeholder.

scripts
-------
- `test-tools-list.sh`: sends `tools/list` to the server using `go run` with configurable `PROMPTS_DIR`.
- `test-prompts.sh`: sends `tools/list`, `list_prompts`, and `expand_prompt` requests; defaults to `examples/prompts`.
