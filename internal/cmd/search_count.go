package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/lxkrmr/godoorpc"
)

const searchCountHelp = `Count records matching a domain for an Odoo model.

Usage:
  gindoo [connection flags] search_count <model> [flags]

Arguments:
  model     Technical model name (e.g. res.partner)

Flags:
  --domain    Odoo domain filter (e.g. "[('is_company', '=', True)]")

Examples:
  gindoo search_count res.partner
  gindoo search_count --domain "[('is_company', '=', True)]" res.partner`

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

	var input searchCountInput
	fs.StringVar(&input.domain, "domain", "", "Odoo domain filter")

	if err := fs.Parse(args); err != nil {
		return searchCountInput{}, err
	}

	positional := fs.Args()
	if len(positional) == 0 {
		return searchCountInput{}, fmt.Errorf("model name is required — run 'gindoo search_count --help'")
	}

	input.model = positional[0]
	return input, nil
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
func RunSearchCount(args []string, conn ConnFlags) {
	input, err := parseSearchCountArgs(args)
	if err == flag.ErrHelp {
		os.Exit(0)
	}
	if err != nil {
		write(errorPayload("search_count", err))
		os.Exit(1)
	}

	var parsedDomain godoorpc.Domain
	if input.domain != "" {
		parsedDomain, err = godoorpc.ParseDomain(input.domain)
		if err != nil {
			write(errorPayload("search_count", fmt.Errorf("invalid domain: %w", err)))
			os.Exit(1)
		}
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
