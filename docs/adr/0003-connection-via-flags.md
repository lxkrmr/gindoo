# ADR 0003: Connection via flags

## Status

Accepted

## Context

`gindoo` needs to know the Odoo URL, database, user, and password to
connect. Common approaches include config files, environment variables,
and CLI flags.

Environment variables for secrets are a security risk: they are visible
in process listings, inherited by child processes, and can leak into
logs. Config files are the secure option but add complexity — file
location, format, creation, and management all need to be handled.

For a tool used locally against a single Odoo instance, a shell alias
covering the flags is a practical and transparent alternative to a
config system.

## Decision

Connection credentials are passed as CLI flags:
`--url`, `--db`, `--user`, `--password`.

No config file system and no environment variables for credentials.
A shell alias is the recommended way to avoid repeating flags:

    alias gindoo='gindoo --url http://localhost:8069 --db mydb --user admin --password secret'

If multiple environments become a real need, a config file with strict
permissions (`chmod 600`) can be introduced at that point.

## Consequences

Positive:
- no config file to manage
- transparent: credentials are visible at the call site
- secure: no env var leakage

Negative:
- flags must be repeated on every call without an alias
- password appears in shell history unless the alias is set up
