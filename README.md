# gindoo

A read-only CLI for inspecting Odoo data.

## What it is

`gindoo` is a tool for exploring a local Odoo instance together with a
coding assistant. The typical workflow:

1. You tell your coding assistant: _"have a look at gindoo"_
2. The assistant reads `gindoo --help` and understands the tool
3. You provide the connection details: URL, database, user, password
4. You and the assistant explore Odoo together — the assistant runs
   `gindoo` commands, you ask questions

`gindoo` is read-only by design. It is safe to run against any Odoo
instance without risk of accidental data mutation.

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

## Usage

Connection flags are required for every command and must come before
the command name:

```sh
gindoo --url <url> --db <db> --user <user> --password <password> <command> [args]
```

### Commands

```sh
# search and read records
gindoo search_read res.partner "[]" "['name', 'email']"
gindoo search_read res.partner "[('is_company', '=', True)]" "['name', 'email']" --limit 5

# count records
gindoo search_count res.partner "[]"
gindoo search_count res.partner "[('is_company', '=', True)]"

# inspect fields of a model
gindoo fields_get res.partner
gindoo fields_get res.partner "['name', 'email']"
```

All output is JSON.

Run `gindoo <command> --help` for command-specific usage and examples.

## License

MIT
