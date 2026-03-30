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

Every response is a JSON object with at minimum an `ok` boolean field.
Successful responses include the requested data. Error responses include
a message field.

## Consequences

Positive:
- single output format to maintain
- output is always machine-readable and pipeable
- consistent structure across all commands

Negative:
- raw JSON is noisier than formatted text for casual terminal use
- requires `jq` or similar for comfortable human reading
