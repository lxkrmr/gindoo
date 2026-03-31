package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/lxkrmr/godoorpc"
)

const searchReadHelp = `Search and read records for an Odoo model.

Usage:
  gindoo [connection flags] search_read <model> <domain> <fields> [--limit N]

Arguments:
  model     Technical model name (e.g. res.partner)
  domain    Odoo domain filter in Odoo list syntax
            Use "[]" for all records, or e.g. "[('is_company', '=', True)]"
  fields    Fields to return in Odoo list syntax, e.g. "['id', 'name']"

Flags:
  --limit   Maximum number of records to return (default: 10)

Examples:
  gindoo --url http://localhost:8069 --db mydb --user admin --password secret search_read res.partner "[]" "['name', 'email']"
  gindoo --url http://localhost:8069 --db mydb --user admin --password secret search_read res.partner "[('is_company', '=', True)]" "['name', 'email']" --limit 5
  gindoo --url http://localhost:8069 --db mydb --user admin --password secret search_read res.partner "[|, ('name', 'ilike', 'foo'), ('id', 'in', [1,2,3])]" "['name', 'email']"`

// searchReadInput holds the parsed data for a search_read command.
type searchReadInput struct {
	model  string
	domain string
	fields string
	limit  int
}

// parseSearchReadArgs parses flags and positional args — calculation.
func parseSearchReadArgs(args []string) (searchReadInput, error) {
	fs := flag.NewFlagSet("search_read", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	fs.Usage = func() { fmt.Println(searchReadHelp) }

	var input searchReadInput
	fs.IntVar(&input.limit, "limit", 10, "Maximum number of records")

	if err := fs.Parse(hoistFlags(args)); err != nil {
		return searchReadInput{}, err
	}

	positional := fs.Args()
	if len(positional) < 3 {
		return searchReadInput{}, fmt.Errorf("model, domain, and fields are required — run 'gindoo search_read --help'")
	}
	if len(positional) > 3 {
		return searchReadInput{}, fmt.Errorf(
			"unexpected argument %q\n"+
				"search_read takes exactly: <model> <domain> <fields> [--limit N]\n"+
				"run 'gindoo search_read --help' for usage",
			positional[3],
		)
	}

	input.model = positional[0]
	input.domain = positional[1]
	input.fields = positional[2]

	return input, nil
}

// buildSearchReadResult shapes the data for the JSON response — pure calculation.
func buildSearchReadResult(input searchReadInput, records any) map[string]any {
	return map[string]any{
		"model":   input.model,
		"domain":  input.domain,
		"fields":  input.fields,
		"limit":   input.limit,
		"records": records,
	}
}

// RunSearchRead executes the search_read command: searches and reads records from an Odoo model.
func RunSearchRead(args []string, conn ConnFlags) {
	input, err := parseSearchReadArgs(args)
	if err == flag.ErrHelp {
		os.Exit(0)
	}
	if err != nil {
		write(errorPayload("search_read", err))
		os.Exit(1)
	}

	parsedDomain, err := godoorpc.ParseDomain(input.domain)
	if err != nil {
		write(errorPayload("search_read", fmt.Errorf("invalid domain %q: %w", input.domain, err)))
		os.Exit(1)
	}

	parsedFields, err := parseFieldList(input.fields)
	if err != nil {
		write(errorPayload("search_read", err))
		os.Exit(1)
	}

	client, err := conn.Connect()
	if err != nil {
		write(errorPayload("search_read", fmt.Errorf("cannot connect to Odoo: %w", err)))
		os.Exit(1)
	}

	records, err := client.ExecuteKW(input.model, "search_read",
		godoorpc.Args{parsedDomain},
		godoorpc.KWArgs{
			"fields": parsedFields,
			"limit":  input.limit,
			"order":  "id asc",
		},
	)
	if err != nil {
		write(errorPayload("search_read", fmt.Errorf("search_read failed for model %q: %w", input.model, err)))
		os.Exit(1)
	}

	write(successPayload("search_read", buildSearchReadResult(input, records)))
}
