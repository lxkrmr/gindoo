package main

import (
	"fmt"
	"os"

	"github.com/lxkrmr/gindoo/internal/cmd"
)

const help = `gindoo — read-only CLI for inspecting Odoo data

Usage:
  gindoo <command> [args]

Commands:
  context       Manage connection contexts
  search_read   Search and read records for a model
  search_count  Count records matching a domain
  fields_get    Describe fields and their metadata for a model
  read_group    Group records by fields and return aggregates

Examples:
  gindoo context create mydev
  gindoo context list
  gindoo context use mydev
  gindoo search_read res.partner "[]" "['name', 'email']"
  gindoo search_read res.partner "[('is_company', '=', True)]" "['name', 'email']" --limit 5
  gindoo search_count res.partner "[('is_company', '=', True)]"
  gindoo fields_get res.partner
  gindoo fields_get res.partner "['name', 'email']"
  gindoo read_group product.template "[]" "['fine_weight:avg']" "['default_code']"

Run 'gindoo <command> --help' for command-specific usage.`

func main() {
	if len(os.Args) < 2 {
		fmt.Println(help)
		os.Exit(0)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "context":
		cmd.RunContext(args)
	case "search_read":
		cmd.RunSearchRead(args)
	case "search_count":
		cmd.RunSearchCount(args)
	case "fields_get":
		cmd.RunFieldsGet(args)
	case "read_group":
		cmd.RunReadGroup(args)
	case "help":
		fmt.Println(help)
	default:
		cmd.WriteError("", fmt.Errorf("unknown command %q — run gindoo --help", command))
		os.Exit(1)
	}
}
