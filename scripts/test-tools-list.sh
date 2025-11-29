#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

if [[ -z "${GOMODCACHE:-}" ]]; then
  export GOMODCACHE="$ROOT_DIR/.gomodcache"
  export GOPATH="$ROOT_DIR/.gopath"
  export GOCACHE="$ROOT_DIR/.gocache"
fi

PROMPTS_DIR="${2:-"$ROOT_DIR/examples/prompts"}"

REQ='{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}'

printf '%s\n' "$REQ" | go run ../cmd/prompts --prompts-dir "$PROMPTS_DIR"