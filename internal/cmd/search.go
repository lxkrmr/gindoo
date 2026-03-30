package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/lxkrmr/godoorpc"
)

const searchHelp = `Search and read records for an Odoo model.

Usage:
  gindoo [connection flags] search <model> [fields...] [flags]

Arguments:
  model     Technical model name (e.g. res.partner)
  fields    Fields to include in the result (default: id only)

Flags:
  --domain    Odoo domain filter (e.g. "[('is_company', '=', True)]")
  --limit     Maximum number of records to return (default: 10)
  --offset    Number of records to skip (default: 0)

Examples:
  gindoo search res.partner
  gindoo search res.partner name email
  gindoo search --domain "[('is_company', '=', True)]" res.partner name email
  gindoo search --limit 5 --offset 10 res.partner name`

// searchInput holds the parsed data for a search command.
type searchInput struct {
	model  string
	fields []string
	domain string
	limit  int
	offset int
}

// parseSearchArgs parses flags and positional args — calculation.
func parseSearchArgs(args []string) (searchInput, error) {
	fs := flag.NewFlagSet("search", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	fs.Usage = func() { fmt.Println(searchHelp) }

	var input searchInput
	fs.StringVar(&input.domain, "domain", "", "Odoo domain filter")
	fs.IntVar(&input.limit, "limit", 10, "Maximum number of records")
	fs.IntVar(&input.offset, "offset", 0, "Number of records to skip")

	if err := fs.Parse(args); err != nil {
		return searchInput{}, err
	}

	positional := fs.Args()
	if len(positional) == 0 {
		return searchInput{}, fmt.Errorf("model name is required — run 'gindoo search --help'")
	}

	input.model = positional[0]
	input.fields = positional[1:]
	if len(input.fields) == 0 {
		input.fields = []string{"id"}
	}

	return input, nil
}

// buildSearchResult shapes the data for the JSON response — pure calculation.
func buildSearchResult(input searchInput, records any) map[string]any {
	return map[string]any{
		"model":   input.model,
		"domain":  input.domain,
		"fields":  input.fields,
		"limit":   input.limit,
		"offset":  input.offset,
		"records": records,
	}
}

// RunSearch executes the search command: reads records from an Odoo model.
func RunSearch(args []string, conn ConnFlags) {
	input, err := parseSearchArgs(args)
	if err == flag.ErrHelp {
		os.Exit(0)
	}
	if err != nil {
		write(errorPayload("search", err))
		os.Exit(1)
	}

	var parsedDomain godoorpc.Domain
	if input.domain != "" {
		parsedDomain, err = godoorpc.ParseDomain(input.domain)
		if err != nil {
			write(errorPayload("search", fmt.Errorf("invalid domain: %w", err)))
			os.Exit(1)
		}
	}

	client, err := conn.Connect()
	if err != nil {
		write(errorPayload("search", fmt.Errorf("cannot connect to Odoo: %w", err)))
		os.Exit(1)
	}

	records, err := client.ExecuteKW(input.model, "search_read",
		godoorpc.Args{parsedDomain},
		godoorpc.KWArgs{
			"fields": input.fields,
			"limit":  input.limit,
			"offset": input.offset,
			"order":  "id asc",
		},
	)
	if err != nil {
		write(errorPayload("search", fmt.Errorf("search failed for model %q: %w", input.model, err)))
		os.Exit(1)
	}

	write(successPayload("search", buildSearchResult(input, records)))
}
