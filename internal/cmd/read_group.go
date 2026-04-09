package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/lxkrmr/godoorpc"
)

const readGroupHelp = `Group records by fields and return aggregates.

Usage:
  gindoo read_group <model> <domain> <fields> <groupby> [--limit N] [--orderby "field"]

Arguments:
  model     Technical model name (e.g. product.template)
  domain    Odoo domain filter (e.g. "[]")
  fields    Fields to aggregate (e.g. "['fine_weight:avg']")
  groupby   Fields to group by (e.g. "['default_code']")

Flags:
  --limit    Maximum number of groups to return (default: 10)
  --orderby  Field to sort groups (e.g. "default_code desc")

Examples:
  gindoo read_group product.template "[]" "['fine_weight:avg']" "['default_code']"
  gindoo read_group res.partner "[('is_company', '=', True)]" "['id:count']" "['country_id']" --limit 5

Uses the current context. Set it with: gindoo context use <name>`

// readGroupInput holds the parsed data for a read_group command.
type readGroupInput struct {
	model   string
	domain  string
	fields  string
	groupby string
	limit   int
	orderby string
}

// parseReadGroupArgs parses flags and positional args — calculation.
func parseReadGroupArgs(args []string) (readGroupInput, error) {
	fs := flag.NewFlagSet("read_group", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	fs.Usage = func() { fmt.Println(readGroupHelp) }

	var input readGroupInput
	fs.IntVar(&input.limit, "limit", 10, "Maximum number of groups")
	fs.StringVar(&input.orderby, "orderby", "", "Field to sort groups")

	if err := fs.Parse(hoistFlags(args)); err != nil {
		return readGroupInput{}, err
	}

	positional := fs.Args()
	if len(positional) < 4 {
		return readGroupInput{}, fmt.Errorf("model, domain, fields, and groupby are required")
	}
	if len(positional) > 4 {
		return readGroupInput{}, fmt.Errorf("unexpected argument %q", positional[4])
	}

	input.model = positional[0]
	input.domain = positional[1]
	input.fields = positional[2]
	input.groupby = positional[3]

	return input, nil
}

// buildReadGroupResult shapes the data for the JSON response — pure calculation.
func buildReadGroupResult(input readGroupInput, records any) map[string]any {
	return map[string]any{
		"model":   input.model,
		"domain":  input.domain,
		"fields":  input.fields,
		"groupby": input.groupby,
		"limit":   input.limit,
		"orderby": input.orderby,
		"records": records,
	}
}

// RunReadGroup executes the read_group command: groups records by fields and returns aggregates.
func RunReadGroup(args []string) {
	input, err := parseReadGroupArgs(args)
	if err == flag.ErrHelp {
		os.Exit(0)
	}
	if err != nil {
		write(errorPayload("read_group", err))
		os.Exit(1)
	}

	_, ctx, err := GetCurrentContext()
	if err != nil {
		write(errorPayload("read_group", err))
		os.Exit(1)
	}

	conn := ConvertContextToConnFlags(ctx)

	parsedDomain, err := godoorpc.ParseDomain(input.domain)
	if err != nil {
		write(errorPayload("read_group", fmt.Errorf("invalid domain %q: %w", input.domain, err)))
		os.Exit(1)
	}

	parsedFields, err := parseFieldList(input.fields)
	if err != nil {
		write(errorPayload("read_group", err))
		os.Exit(1)
	}

	parsedGroupBy, err := parseFieldList(input.groupby)
	if err != nil {
		write(errorPayload("read_group", err))
		os.Exit(1)
	}

	client, err := conn.Connect()
	if err != nil {
		write(errorPayload("read_group", fmt.Errorf("cannot connect to Odoo: %w", err)))
		os.Exit(1)
	}

	// Call Odoo's read_group
	records, err := client.ExecuteKW(input.model, "read_group",
		godoorpc.Args{parsedDomain, parsedFields, parsedGroupBy},
		godoorpc.KWArgs{
			"limit":   input.limit,
			"orderby": input.orderby,
			"lazy":    false, // Force DB query
		},
	)
	if err != nil {
		write(errorPayload("read_group", fmt.Errorf("read_group failed: %w", err)))
		os.Exit(1)
	}

	write(successPayload("read_group", buildReadGroupResult(input, records)))
}