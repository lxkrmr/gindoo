# ADR 0005: Help-driven discovery

## Status

Accepted

## Context

`gindoo` should be usable by a developer or a coding agent who has never
read the README or any documentation. The only guaranteed entry point is
the CLI itself.

This means `--help` must carry enough information to answer:
- what is gindoo?
- what can I do with it?
- how do I connect to Odoo?
- what does this specific command do, and how do I call it?

Error messages must answer:
- what went wrong?
- what should I try next?

## Decision

`gindoo --help` describes the tool's purpose and lists all commands with
a one-line description each.

Every subcommand help (`gindoo search --help`) includes:
- what the command does
- all flags with descriptions and defaults
- at least one concrete usage example

Error messages always include a human-readable explanation and, where
possible, a concrete suggestion for what to do next.

No separate `describe` or `about` command is added. If `--help` is not
sufficient for discovery, the help text is improved — not a new command.

## Consequences

Positive:
- a human or agent can fully discover gindoo from the CLI alone
- no dependency on README or external documentation
- clear error messages reduce debugging time

Negative:
- help text requires care to stay accurate and useful as the tool evolves
