Decisions
=========

Prompt Storage
--------------
- Store prompts as individual `.yaml` files in a single directory to keep authoring simple and editable without code changes.

Directory Validation
--------------------
- Require an absolute `--prompts-dir`, reject root `/`, and resolve symlinks to avoid accidental traversal or invalid inputs when launching the server.

Placeholder Handling
--------------------
- When `{{input}}` is present in a prompt, replace all occurrences; otherwise append the input after two newlines to preserve templates that do not declare placeholders explicitly.

Logging
-------
- Use zap production JSON logging by default for structured output suitable for tooling; defer `logger.Sync()` in main to flush buffers.
