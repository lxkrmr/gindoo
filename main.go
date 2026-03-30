package main

import (
	"fmt"
	"os"

	"github.com/lxkrmr/gindoo/cmd"
)

const help = `gindoo — read-only CLI for inspecting Odoo data

Usage:
  gindoo [connection flags] <command> [args]

Commands:
  search      Search and read records for a model
  read        Read fields from a single record by ID
  fields_get  Describe fields and their metadata for a model

Connection flags (required for all commands):
  --url       Odoo base URL (e.g. http://localhost:8069)
  --db        Database name
  --user      Login user
  --password  Login password

Examples:
  gindoo search res.partner name email --url http://localhost:8069 --db mydb --user admin --password secret
  gindoo read res.partner 1 name email --url http://localhost:8069 --db mydb --user admin --password secret
  gindoo fields_get res.partner --url http://localhost:8069 --db mydb --user admin --password secret

Tip: use a shell alias to avoid repeating connection flags:
  alias gindoo='gindoo --url http://localhost:8069 --db mydb --user admin --password secret'

Run 'gindoo <command> --help' for command-specific usage.`

func main() {
	if len(os.Args) < 2 {
		fmt.Println(help)
		os.Exit(0)
	}

	switch os.Args[1] {
	case "search":
		cmd.RunSearch(os.Args[2:])
	case "read":
		cmd.RunRead(os.Args[2:])
	case "fields_get":
		cmd.RunFieldsGet(os.Args[2:])
	case "--help", "-h", "help":
		fmt.Println(help)
	default:
		fmt.Fprintf(os.Stdout,
			`{"ok":false,"command":"","error":"unknown command %q — run gindoo --help"}`,
			os.Args[1],
		)
		fmt.Println()
		os.Exit(1)
	}
}
