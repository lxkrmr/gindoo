# gindoo

A read-only CLI for inspecting Odoo data.

`gindoo` is designed for two kinds of users working together: a developer
who knows what they want to find out, and a coding assistant who uses
gindoo to explore the Odoo instance and surface the answer.

It gives both a safe window into a local Odoo instance — more detail
than the UI, more context than the database, no risk of accidental
mutation.

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

## Quickstart

1. Add an alias to your `~/.zshrc` or `~/.bashrc`:

```sh
alias gindoo='gindoo --url http://localhost:8069 --db mydb --user admin --password secret'
```

2. Reload your shell:

```sh
source ~/.zshrc  # or ~/.bashrc
```

3. Start inspecting:

```sh
gindoo search res.partner name email
gindoo search_count product.product
gindoo fields_get sale.order
```

## Usage

Connection flags are required for every command and must come before
the command name:

```sh
gindoo --url http://localhost:8069 --db mydb --user admin --password secret <command>
```

With the alias from the quickstart, this becomes simply:

```sh
gindoo <command>
```

### Commands

```sh
# search records
gindoo search res.partner name email
gindoo search --domain "[('is_company', '=', True)]" --limit 5 res.partner name

# count records
gindoo search_count res.partner
gindoo search_count --domain "[('is_company', '=', True)]" res.partner

# read a single record by ID
gindoo read res.partner 1 name email phone

# inspect fields of a model
gindoo fields_get res.partner
gindoo fields_get res.partner name email
```

All output is JSON.

## License

MIT
