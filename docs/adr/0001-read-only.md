# ADR 0001: Read-only by design

## Status

Accepted

## Context

`gindoo` is an inspection tool for Odoo developers. Its purpose is to
give developers a safe, detailed view into what is actually in a local
Odoo instance — field types, record values, domain-based searches.

Adding mutation commands (write, create, unlink) would create a path
to change Odoo data that bypasses the UI and the ORM-level business
logic the UI exercises. A wrong write via CLI produces data that may
go unnoticed until it surfaces in the UI. Unlike a bad read, a bad
write cannot be trivially ignored.

The goal of Odoo development is to make things manageable through the
UI. `gindoo` exists to help verify that — not to change it.

## Decision

`gindoo` is read-only. Mutation commands are out of scope.

The following will not be added without a concrete justification:
- `write`
- `create`
- `unlink`
- any general mutation primitive

## Consequences

Positive:
- no risk of accidental data mutation
- clear, focused scope
- safe to run against any Odoo instance without side effects

Negative:
- developers who need to mutate data must use the UI or tested
  application code
