# ADR 0006: internal/cmd instead of cmd

## Status

Accepted

## Context

`gindoo` is a CLI tool, not a library. The `cmd` package contains the
command implementations and is only ever imported by `main.go`. It is
not part of any public API.

Placing it at `cmd/` makes it importable by any external Go module,
which signals that it is a public API — which it is not.

## Decision

The command package lives at `internal/cmd/`.

Go's `internal/` convention prevents any module outside of
`github.com/lxkrmr/gindoo` from importing it. This makes the
intent clear: this code is private to the tool.

## Consequences

Positive:
- clearly signals that cmd is not a public API
- prevents accidental external imports

Negative:
- none for a CLI tool
