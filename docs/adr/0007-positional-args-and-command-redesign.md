# ADR 0007: Positional arguments and command redesign

## Status

Accepted

## Context

The original `search` command had a positional `model`, variadics `fields`,
and optional flags `--domain`, `--limit`, `--offset`. This caused a
silent failure when `--domain` was placed after the fields: Go's `flag`
package stops parsing at the first non-flag argument, so `--domain` was
treated as a field name and sent to Odoo, which returned a server error.

The root cause was not the flag parser — it was a design mismatch:
required arguments were flags, optional arguments were variadics, and
the two were mixed in a way that made position matter without saying so.

Reviewing the design from first principles revealed a cleaner approach:

- Required arguments belong in positional positions.
- Optional arguments belong in flags.
- Variadics mixed with optional flags always create ordering ambiguity.

Additionally, `search` was a misleading name: it calls Odoo's
`search_read` method internally. The command name was hiding which RPC
method was being called, which contradicts ADR 0004.

The `read` command was also found to be redundant: any call to `read`
can be expressed as `search_read` with a domain of `[('id', '=', X)]`.
Removing it reduces surface area without removing capability.

## Decision

### `search` is renamed to `search_read`

The command name now matches the Odoo RPC method it calls, consistent
with ADR 0004.

### `domain` is a required positional argument

`--domain` as an optional flag is replaced by `<domain>` as the second
positional argument. A search without a domain is conceptually
incomplete — even an empty domain `[]` is an explicit choice.

### `fields` is a required positional argument in Odoo list syntax

`fields` as variadics is replaced by `<fields>` as the third positional
argument. The format follows Odoo's own syntax: `"['name', 'email']"`.
This is consistent with the domain format and with how Odoo's own API
accepts field lists.

`fields_get` adopts the same format for its optional fields argument:
`gindoo fields_get res.partner "['name', 'email']"`.

### `--limit` remains an optional flag with default 10

`limit` is a modifier — not a subject of the command. In Odoo's own
API it is a keyword argument. `--limit 5` is self-documenting and
consistent with standard CLI conventions.

`--offset` is removed. Pagination via offset is not a primary use case
for an inspection tool. It can be reintroduced if a concrete need arises.

### `read` is removed

`search_read` with an ID domain covers the same use case. The `read`
command is redundant.

### Usage lines and examples must be consistent

Usage lines and examples in every help text must show the same argument
format. Examples include the full command with connection flags, because
that is what an agent constructs.

## Resulting command signatures

```
gindoo search_read <model> <domain> <fields> [--limit N]
gindoo search_count <model> <domain>
gindoo fields_get <model> [fields]
```

Examples:

```
gindoo --url http://localhost:8069 --db mydb --user admin --password secret \
  search_read res.partner "[]" "['name', 'email']"

gindoo --url http://localhost:8069 --db mydb --user admin --password secret \
  search_read res.partner "[('is_company', '=', True)]" "['name', 'email']" --limit 5

gindoo --url http://localhost:8069 --db mydb --user admin --password secret \
  search_count res.partner "[('is_company', '=', True)]"

gindoo --url http://localhost:8069 --db mydb --user admin --password secret \
  fields_get res.partner

gindoo --url http://localhost:8069 --db mydb --user admin --password secret \
  fields_get res.partner "['name', 'email']"
```

## Consequences

Positive:
- no flag ordering ambiguity
- required arguments are always present — no silent empty-domain searches
- command names match Odoo RPC method names (ADR 0004)
- usage lines and examples are consistent (ADR 0005)
- smaller surface area: one fewer command, one fewer flag

Negative:
- breaking change for existing callers
- `fields` as a quoted list is slightly more verbose than space-separated
  variadics for short field lists
