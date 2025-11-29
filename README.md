prompts-mcp
=============
Small Go MCP server that reads prompt templates from YAML files and exposes them over stdio using JSON-RPC 2.0 tools.

Overview
--------
- Serves two MCP tools: `list_prompts` (lists YAML prompt files) and `expand_prompt` (returns the prompt with `{{input}}` replaced, or appends input if no placeholder).
- Accepts an absolute `--prompts-dir` pointing to the directory containing `.yaml` prompt definitions.
- Ships example prompts and bash scripts to exercise the MCP server over stdio.

Project Structure
-----------------
- `cmd/prompts/` — CLI entrypoint; wires flags, logging, validation, server.
- `internal/` — core packages:
  - `logging/` — zap production logger setup.
  - `mcp/` — JSON-RPC handling, MCP tools, stdio server.
  - `prompts/` — prompt model and YAML-backed repository.
  - `validate/` — path validation helpers and tests.
- `examples/prompts/` — sample prompt YAML files.
- `scripts/` — helper scripts to list/call tools against the server.
- `Makefile` — build/test/vet targets.

Installation & Setup
--------------------
- Prereqs: Go 1.23+.
- Install dependencies: `go mod download`
- Build: `go build -o prompts ./cmd/prompts` (or `make build`)
- Test: `go test ./...` (or `make test`)
- Static analysis: `make vet`
- Ensure you have a prompts directory with `.yaml` files and pass its absolute path to `--prompts-dir`.

Usage
-----
- Start the server over stdio, supplying an absolute prompts directory:
  - `go run ./cmd/prompts --prompts-dir /absolute/path/to/prompts`
- Example JSON-RPC interactions (one per line):
  - `{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}`
  - `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"list_prompts","arguments":{}}}`
  - `{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"expand_prompt","arguments":{"command":"research","input":"Example topic"}}}`
- Helper scripts (set `PROMPTS_DIR` to override default `examples/prompts`):
  - `./scripts/test-tools-list.sh`
  - `./scripts/test-prompts.sh`
  - Scripts set `GOMODCACHE`, `GOPATH`, and `GOCACHE` under the repo when those env vars are unset.

Architecture
------------
- See `docs/architecture.md` for data flow from CLI flag parsing through prompt loading to MCP responses.

Documentation
-------------
- `docs/architecture.md`
- `docs/api.md`
- `docs/modules.md`
- `docs/setup.md`
- `docs/research/`
- `docs/decisions.md`

Tech Stack
----------
- Go 1.23
- JSON-RPC 2.0 over stdio (Model Context Protocol 2024-11-05)
- zap for logging
- YAML (gopkg.in/yaml.v3) for prompt storage
