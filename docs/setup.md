Setup
=====

Prerequisites
-------------
- Go 1.23 or newer.
- A directory of prompt files ending in `.md`.

Build and Test
--------------
- Install deps: `go mod download`
- Build binary: `go build -o prompts ./cmd/prompts` (or `make build`)
- Run tests: `go test ./...` (or `make test`)
- Static analysis: `make vet`

Running the Server
------------------
- The prompts directory flag must be absolute:
  - `go run ./cmd/prompts --prompts-dir /absolute/path/to/prompts`
- The server reads JSON-RPC 2.0 requests from stdin and writes responses to stdout; log output goes to stderr (zap production JSON).

Local Cache Configuration
-------------------------
- Helper scripts set per-repo caches when unset:
  - `GOMODCACHE=$ROOT/.gomodcache`
  - `GOPATH=$ROOT/.gopath`
  - `GOCACHE=$ROOT/.gocache`

Prompt File Format
------------------
- Each prompt file must be Markdown with extension `.md` and include:
  - YAML frontmatter wrapped by `---` delimiters.
  - `description: "short description"` in frontmatter.
  - Prompt template body in the Markdown content after frontmatter.
- Optional placeholder: `{{input}}` will be replaced with the `input` argument when expanding. If absent, the input is appended after two newlines.
- Example: see `examples/prompts/research.md` and `examples/prompts/summarize.md`.

Quick Checks
------------
- Run from the `scripts/` directory (the scripts use `go run ../cmd/prompts`):
  - `cd scripts && ./test-tools-list.sh`
  - `cd scripts && ./test-prompts.sh`
- Optional script overrides:
  - `./test-tools-list.sh _ /absolute/path/to/prompts`
  - `./test-prompts.sh _ /absolute/path/to/prompts summarize "text to summarize"`
