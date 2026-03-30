# gindoo

A read-only CLI for inspecting Odoo data.

`gindoo` gives developers a safe window into a local Odoo instance —
more detail than the UI, more context than the database.

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

Connection flags are required for every command:

```sh
gindoo --url http://localhost:8069 --db mydb --user admin --password secret <command>
```

Set up an alias to avoid repeating them:

```sh
alias gindoo='gindoo --url http://localhost:8069 --db mydb --user admin --password secret'
```

### Commands

```sh
# search records
gindoo search res.partner name email
gindoo search res.partner name --domain "[('is_company', '=', True)]" --limit 5

# read a single record by ID
gindoo read res.partner 1 name email phone

# inspect fields of a model
gindoo fields_get res.partner
gindoo fields_get res.partner name email
```

All output is JSON.

## License

MIT
