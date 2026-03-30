package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/lxkrmr/godoorpc"
)

const fieldsGetHelp = `Describe fields and their metadata for an Odoo model.

Usage:
  gindoo fields_get [flags] <model> [fields...]

Arguments:
  model     Technical model name (e.g. res.partner)
  fields    Specific field names to inspect (default: all fields)

Connection flags:
  --url       Odoo base URL (e.g. http://localhost:8069)
  --db        Database name
  --user      Login user
  --password  Login password

Examples:
  gindoo fields_get res.partner
  gindoo fields_get res.partner name email phone
  gindoo fields_get sale.order state amount_total`

// fieldsGetInput holds the parsed data for a fields_get command.
type fieldsGetInput struct {
	conn   connFlags
	model  string
	fields []string
}

// parseFieldsGetArgs parses flags and positional args — calculation.
func parseFieldsGetArgs(args []string) (fieldsGetInput, error) {
	fs := flag.NewFlagSet("fields_get", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	fs.Usage = func() { fmt.Println(fieldsGetHelp) }

	var input fieldsGetInput
	registerConnFlags(fs, &input.conn)

	if err := fs.Parse(args); err != nil {
		return fieldsGetInput{}, err
	}

	positional := fs.Args()
	if len(positional) == 0 {
		return fieldsGetInput{}, fmt.Errorf("model name is required — run 'gindoo fields_get --help'")
	}

	input.model = positional[0]
	input.fields = positional[1:]

	return input, nil
}

// fieldArgs builds the first argument for fields_get — pure calculation.
// Odoo expects a list of field names, or false to get all fields.
func fieldArgs(fields []string) any {
	if len(fields) > 0 {
		return fields
	}
	return false
}

// buildFieldsGetResult shapes the data for the JSON response — pure calculation.
func buildFieldsGetResult(input fieldsGetInput, fields any) map[string]any {
	return map[string]any{
		"model":  input.model,
		"fields": fields,
	}
}

// RunFieldsGet orchestrates side effects: parse, connect, execute, write.
func RunFieldsGet(args []string) {
	input, err := parseFieldsGetArgs(args)
	if err == flag.ErrHelp {
		os.Exit(0)
	}
	if err != nil {
		write(errorPayload("fields_get", err))
		os.Exit(1)
	}

	client, err := input.conn.connect()
	if err != nil {
		write(errorPayload("fields_get", fmt.Errorf("cannot connect to Odoo: %w", err)))
		os.Exit(1)
	}

	result, err := client.ExecuteKW(input.model, "fields_get",
		godoorpc.Args{fieldArgs(input.fields)},
		godoorpc.KWArgs{
			"attributes": []string{"string", "type", "required", "readonly", "relation", "selection"},
		},
	)
	if err != nil {
		write(errorPayload("fields_get", fmt.Errorf("fields_get failed for model %q: %w", input.model, err)))
		os.Exit(1)
	}

	write(successPayload("fields_get", buildFieldsGetResult(input, result)))
}
