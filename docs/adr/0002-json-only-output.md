# ADR 0002: JSON-only output

## Status

Accepted

## Context

CLI tools typically offer human-readable text output by default and
machine-readable output as an option. `gindoo` is primarily used by
developers inspecting Odoo data — often piping results into other tools,
agents, or scripts.

Supporting two output modes (text and JSON) doubles the surface area to
maintain and introduces the question of which mode is the default. Text
output that works well for humans tends to lose structure that machines
need.

## Decision

`gindoo` always outputs JSON. There is no text mode and no `--output`
flag.

Every response follows this structure:

```json
{
  "ok": true,
  "command": "search",
  "data": { ... }
}
```

On error:

```json
{
  "ok": false,
  "command": "search",
  "error": "cannot connect to Odoo at localhost:8069"
}
```

- `ok` — boolean, always present. Enables reliable success/failure
  detection without parsing error messages.
- `command` — the name of the command that was called.
- `data` — present on success. Contains the command-specific payload.
- `error` — present on failure. A human-readable error message.

`next_commands` is intentionally omitted. Unlike a workflow tool,
gindoo's commands are independent — what to inspect next depends on
what the developer or agent sees in the data, not on a fixed sequence.

## Consequences

Positive:
- single output format to maintain
- output is always machine-readable and pipeable
- consistent structure across all commands
- `ok` field enables simple success checks in scripts and agents

Negative:
- raw JSON is noisier than formatted text for casual terminal use
- requires `jq` or similar for comfortable human reading
