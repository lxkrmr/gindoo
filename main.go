package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lxkrmr/gindoo/internal/cmd"
)

const help = `gindoo — read-only CLI for inspecting Odoo data

Usage:
  gindoo --url <url> --db <db> --user <user> --password <password> <command> [args]

Commands:
  search_read   Search and read records for a model
  search_count  Count records matching a domain
  fields_get    Describe fields and their metadata for a model

Connection flags (required, must come before the command):
  --url       Odoo base URL (e.g. http://localhost:8069)
  --db        Database name
  --user      Login user
  --password  Login password

Examples:
  gindoo --url http://localhost:8069 --db mydb --user admin --password secret search_read res.partner "[]" "['name', 'email']"
  gindoo --url http://localhost:8069 --db mydb --user admin --password secret search_read res.partner "[('is_company', '=', True)]" "['name', 'email']" --limit 5
  gindoo --url http://localhost:8069 --db mydb --user admin --password secret search_count res.partner "[('is_company', '=', True)]"
  gindoo --url http://localhost:8069 --db mydb --user admin --password secret fields_get res.partner
  gindoo --url http://localhost:8069 --db mydb --user admin --password secret fields_get res.partner "['name', 'email']"

Run 'gindoo <command> --help' for command-specific usage.`

func main() {
	fs := flag.NewFlagSet("gindoo", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	fs.Usage = func() { fmt.Println(help) }

	var conn cmd.ConnFlags
	cmd.RegisterConnFlags(fs, &conn)

	if err := fs.Parse(os.Args[1:]); err != nil {
		if err == flag.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}

	remaining := fs.Args()
	if len(remaining) == 0 {
		fmt.Println(help)
		os.Exit(0)
	}

	switch remaining[0] {
	case "search_read":
		cmd.RunSearchRead(remaining[1:], conn)
	case "search_count":
		cmd.RunSearchCount(remaining[1:], conn)
	case "fields_get":
		cmd.RunFieldsGet(remaining[1:], conn)
	case "help":
		fmt.Println(help)
	default:
		cmd.WriteError("", fmt.Errorf("unknown command %q — run gindoo --help", remaining[0]))
		os.Exit(1)
	}
}
