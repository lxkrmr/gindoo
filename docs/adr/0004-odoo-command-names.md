# ADR 0004: Command names follow Odoo conventions

## Status

Accepted

## Context

`gindoo` wraps Odoo RPC calls. The three operations it exposes map
directly to Odoo model methods: `search_read`, `read`, and `fields_get`.

Using the same names as Odoo reduces the mental translation for
developers who already know the Odoo ORM. It also makes clear which
RPC method is being called under the hood.

## Decision

`gindoo` commands are named after the Odoo methods they call:

- `gindoo search` → `search_read`
- `gindoo read` → `read`
- `gindoo fields_get` → `fields_get`

## Consequences

Positive:
- zero translation for developers who know Odoo
- makes the underlying RPC call obvious
- consistent with how the original Python indoo tool named its commands

Negative:
- `fields_get` is not a typical CLI command name style
