package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/lxkrmr/godoorpc"
)

const searchCountHelp = `Count records matching a domain for an Odoo model.

Usage:
  gindoo search_count <model> <domain>

Arguments:
  model     Technical model name (e.g. res.partner)
  domain    Odoo domain filter in Odoo list syntax
            Use "[]" for all records, or e.g. "[('is_company', '=', True)]"

Examples:
  gindoo search_count res.partner "[]"
  gindoo search_count res.partner "[('is_company', '=', True)]"

Uses the current context. Set it with: gindoo context use <name>`

// searchCountInput holds the parsed data for a search_count command.
type searchCountInput struct {
	model  string
	domain string
}

// parseSearchCountArgs parses flags and positional args — calculation.
func parseSearchCountArgs(args []string) (searchCountInput, error) {
	fs := flag.NewFlagSet("search_count", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	fs.Usage = func() { fmt.Println(searchCountHelp) }

	if err := fs.Parse(hoistFlags(args)); err != nil {
		return searchCountInput{}, err
	}

	positional := fs.Args()
	if len(positional) < 2 {
		return searchCountInput{}, fmt.Errorf("model and domain are required — run 'gindoo search_count --help'")
	}
	if len(positional) > 2 {
		return searchCountInput{}, fmt.Errorf(
			"unexpected argument %q\n"+
				"search_count takes exactly: <model> <domain>\n"+
				"run 'gindoo search_count --help' for usage",
			positional[2],
		)
	}

	return searchCountInput{
		model:  positional[0],
		domain: positional[1],
	}, nil
}

// buildSearchCountResult shapes the data for the JSON response — pure calculation.
func buildSearchCountResult(input searchCountInput, count any) map[string]any {
	return map[string]any{
		"model":  input.model,
		"domain": input.domain,
		"count":  count,
	}
}

// RunSearchCount executes the search_count command: counts records matching a domain in an Odoo model.
func RunSearchCount(args []string) {
	input, err := parseSearchCountArgs(args)
	if err == flag.ErrHelp {
		os.Exit(0)
	}
	if err != nil {
		write(errorPayload("search_count", err))
		os.Exit(1)
	}

	_, ctx, err := GetCurrentContext()
	if err != nil {
		write(errorPayload("search_count", err))
		os.Exit(1)
	}

	conn := ConvertContextToConnFlags(ctx)

	parsedDomain, err := godoorpc.ParseDomain(input.domain)
	if err != nil {
		write(errorPayload("search_count", fmt.Errorf("invalid domain %q: %w", input.domain, err)))
		os.Exit(1)
	}

	client, err := conn.Connect()
	if err != nil {
		write(errorPayload("search_count", fmt.Errorf("cannot connect to Odoo: %w", err)))
		os.Exit(1)
	}

	count, err := client.ExecuteKW(input.model, "search_count",
		godoorpc.Args{parsedDomain},
		godoorpc.KWArgs{},
	)
	if err != nil {
		write(errorPayload("search_count", fmt.Errorf("search_count failed for model %q: %w", input.model, err)))
		os.Exit(1)
	}

	write(successPayload("search_count", buildSearchCountResult(input, count)))
}
