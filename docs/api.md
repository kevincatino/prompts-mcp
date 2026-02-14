API
===

Protocol
--------
- Transport: JSON-RPC 2.0 over stdio (newline-delimited messages).
- MCP protocol version: `2024-11-05`.
- Methods implemented:
  - `initialize`
  - `notifications/initialized` (ack, no response)
  - `tools/list`
  - `tools/call` (`list_prompts`, `expand_prompt`)

Initialize
----------
- Request:
  - `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","clientInfo":{"name":"example","version":"dev"}}}`
- Response:
  - `{"jsonrpc":"2.0","id":1,"result":{"protocolVersion":"2024-11-05","capabilities":{"tools":{}},"serverInfo":{"name":"codex-subagents","version":"0.1.0"},"clientInfo":{"name":"example","version":"dev"}}}`
- Next: client sends `{"jsonrpc":"2.0","method":"notifications/initialized"}` (ignored by server).

tools/list
----------
- Request:
  - `{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}`
- Response:
  - `{"jsonrpc":"2.0","id":2,"result":{"tools":[{"name":"list_prompts","description":"List available prompt commands with descriptions. Users invoke prompt commands in messages by writing :<command> (e.g., :research topic); this lists what commands exist.","inputSchema":{"type":"object","properties":{}}},{"name":"expand_prompt","description":"Expands a named prompt command into a full prompt. When ANY user message contains a :-prefixed command token (e.g., :research topic, :summarize this article, :critique \"draft\"), you MUST call this tool before responding. Strip the leading colon from the command name and pass any remaining text (including quotes) as the input. Do not bypass this tool or answer directly whenever a :-prefixed command is present.","inputSchema":{"type":"object","properties":{"command":{"type":"string","description":"Prompt command name (filename without extension)"},"input":{"type":"string","description":"Text to inject into the prompt template"}},"required":["command","input"]}}]}}`

tools/call: list_prompts
------------------------
- Request:
  - `{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"list_prompts","arguments":{}}}`
- Success result:
  - `{"jsonrpc":"2.0","id":3,"result":{"content":[{"type":"text","text":"{\"prompts\":[{\"name\":\"research\",\"description\":\"Research a topic and provide a concise summary with key sources.\"},{\"name\":\"summarize\",\"description\":\"Summarize provided text into a tight digest with bullets.\"}]}"}],"structuredContent":{"prompts":[{"name":"research","description":"Research a topic and provide a concise summary with key sources."},{"name":"summarize","description":"Summarize provided text into a tight digest with bullets."}]}}}`
- Errors: surfaced as MCP tool error (`result.isError:true`) with textual message.

tools/call: expand_prompt
-------------------------
- Request:
  - `{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"expand_prompt","arguments":{"command":"research","input":"Example topic"}}}`
- Success result:
  - `{"jsonrpc":"2.0","id":4,"result":{"content":[{"type":"text","text":"{\"prompt\":\"You are a focused researcher. Investigate the topic below and return:\\n- A 3-5 sentence summary\\n- 3 key findings\\n- Source names or links if mentioned in provided context\\n\\nTopic:\\nExample topic\"}"}],"structuredContent":{"prompt":"You are a focused researcher. Investigate the topic below and return:\n- A 3-5 sentence summary\n- 3 key findings\n- Source names or links if mentioned in provided context\n\nTopic:\nExample topic"}}}`
- Placeholder behavior:
  - If the prompt contains `{{input}}`, all instances are replaced with the provided `input`.
  - If no placeholder exists, the `input` is appended after two newlines to preserve the template body.
- Errors:
  - Unknown prompt file or Markdown/frontmatter parse issues surface as `result.isError:true` with message text.

Error Codes
-----------
- JSON-RPC errors:
  - `-32600` invalid request (e.g., wrong jsonrpc version).
  - `-32601` method/tool not found.
  - `-32602` invalid params.
  - `-32603` internal error (handler failures before tool execution).
- Tool errors:
  - Returned via `result.isError:true` with textual message in `content[0].text`.

Notes
-----
- Messages are newline-delimited; avoid embedding raw newlines inside JSON literals unless escaped.
- The server enforces `jsonrpc:"2.0"` in requests and echoes the same in responses.
