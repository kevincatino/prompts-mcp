#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

if [[ -z "${GOMODCACHE:-}" ]]; then
  export GOMODCACHE="$ROOT_DIR/.gomodcache"
  export GOPATH="$ROOT_DIR/.gopath"
  export GOCACHE="$ROOT_DIR/.gocache"
fi

PROMPTS_DIR="${2:-"$ROOT_DIR/examples/prompts"}"
COMMAND="${3:-"research"}"
INPUT="${4:-"LLM routing frameworks"}"

LIST_REQ='{"jsonrpc":"2.0","id":10,"method":"tools/list","params":{}}'
LIST_PROMPTS_REQ='{"jsonrpc":"2.0","id":11,"method":"tools/call","params":{"name":"list_prompts","arguments":{}}}'
EXPAND_REQ_TEMPLATE='{"jsonrpc":"2.0","id":12,"method":"tools/call","params":{"name":"expand_prompt","arguments":{"command":"%s","input":"%s"}}}'

EXPAND_REQ=$(printf "$EXPAND_REQ_TEMPLATE" "$COMMAND" "$INPUT")

printf '%s\n' "$LIST_REQ" | go run ../cmd/prompts --prompts-dir "$PROMPTS_DIR"
printf '%s\n' "$LIST_PROMPTS_REQ" | go run ../cmd/prompts --prompts-dir "$PROMPTS_DIR"
printf '%s\n' "$EXPAND_REQ" | go run ../cmd/prompts --prompts-dir "$PROMPTS_DIR"