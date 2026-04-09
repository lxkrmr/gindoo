# gindoo

A read-only CLI for inspecting Odoo data.

## What it is

`gindoo` is a tool for exploring a local Odoo instance together with a
coding assistant. The typical workflow:

1. You tell your coding assistant: _"have a look at gindoo"_
2. The assistant reads `gindoo --help` and understands the tool
3. You create a connection context (one time)
4. You and the assistant explore Odoo together — the assistant runs
   `gindoo` commands, you ask questions

`gindoo` is a dev tool built for local Odoo development instances.
It is read-only by design - no risk of accidental data mutation.
Do not use production credentials with gindoo.

## Install

```sh
go install github.com/lxkrmr/gindoo@latest
```

Requires Go. The binary lands in `~/go/bin/gindoo`, which should already
be in your `$PATH` if you have used `go install` before.

If `@latest` resolves to an older version after a new release, bypass
the module proxy cache with:

```sh
GOPROXY=direct go install github.com/lxkrmr/gindoo@latest
```

## Setup

Before using `gindoo`, create a connection context:

```sh
gindoo context create mydev
```

This will prompt for:
- URL (e.g. http://localhost:8069)
- Database name
- Login user
- Password

The context is saved to `~/.config/gindoo/contexts.json` and can be
reused. If you have multiple Odoo instances:

```sh
gindoo context create mydev
gindoo context create staging
gindoo context list
gindoo context use staging   # switch between contexts
```

## Usage

### Manage contexts

```sh
gindoo context create <name>   # Create a new connection context
gindoo context list            # Show all contexts (current marked with *)
gindoo context use <name>      # Set as current context
gindoo context remove <name>   # Delete a context
```

### Query Odoo

```sh
# search and read records
gindoo search_read res.partner "[]" "['name', 'email']"
gindoo search_read res.partner "[('is_company', '=', True)]" "['name', 'email']" --limit 5

# count records
gindoo search_count res.partner "[]"
gindoo search_count res.partner "[('is_company', '=', True)]"

# group and aggregate records
gindoo read_group product.template "[]" "['fine_weight:avg']" "['default_code']"
gindoo read_group res.partner "[('is_company', '=', True)]" "['id:count']" "['country_id']" --limit 5

# inspect fields of a model
gindoo fields_get res.partner
gindoo fields_get res.partner "['name', 'email']"
```

All output is JSON.

Run `gindoo <command> --help` for command-specific usage and examples.

## License

MIT
