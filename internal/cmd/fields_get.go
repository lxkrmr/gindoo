package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/lxkrmr/godoorpc"
)

const fieldsGetHelp = `Describe fields and their metadata for an Odoo model.

Usage:
  gindoo fields_get <model> [fields]

Arguments:
  model     Technical model name (e.g. res.partner)
  fields    Specific fields to inspect in Odoo list syntax (optional, default: all fields)
            e.g. "['name', 'email']"

Examples:
  gindoo fields_get res.partner
  gindoo fields_get res.partner "['name', 'email']"

Uses the current context. Set it with: gindoo context use <name>`

// fieldsGetInput holds the parsed data for a fields_get command.
type fieldsGetInput struct {
	model  string
	fields string // empty means all fields
}

// parseFieldsGetArgs parses flags and positional args — calculation.
func parseFieldsGetArgs(args []string) (fieldsGetInput, error) {
	fs := flag.NewFlagSet("fields_get", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	fs.Usage = func() { fmt.Println(fieldsGetHelp) }

	if err := fs.Parse(hoistFlags(args)); err != nil {
		return fieldsGetInput{}, err
	}

	positional := fs.Args()
	if len(positional) == 0 {
		return fieldsGetInput{}, fmt.Errorf("model name is required — run 'gindoo fields_get --help'")
	}
	if len(positional) > 2 {
		return fieldsGetInput{}, fmt.Errorf(
			"unexpected argument %q\n"+
				"fields_get takes: <model> [fields]\n"+
				"run 'gindoo fields_get --help' for usage",
			positional[2],
		)
	}

	input := fieldsGetInput{model: positional[0]}
	if len(positional) > 1 {
		input.fields = positional[1]
	}
	return input, nil
}

// fieldArgs builds the first argument for fields_get — pure calculation.
// Odoo expects a list of field names, or false to get all fields.
func fieldArgs(fields string) (any, error) {
	if fields == "" {
		return false, nil
	}
	return parseFieldList(fields)
}

// buildFieldsGetResult shapes the data for the JSON response — pure calculation.
func buildFieldsGetResult(input fieldsGetInput, fields any) map[string]any {
	return map[string]any{
		"model":  input.model,
		"fields": fields,
	}
}

// RunFieldsGet executes the fields_get command: describes fields and metadata for an Odoo model.
func RunFieldsGet(args []string) {
	input, err := parseFieldsGetArgs(args)
	if err == flag.ErrHelp {
		os.Exit(0)
	}
	if err != nil {
		write(errorPayload("fields_get", err))
		os.Exit(1)
	}

	_, ctx, err := GetCurrentContext()
	if err != nil {
		write(errorPayload("fields_get", err))
		os.Exit(1)
	}

	conn := ConvertContextToConnFlags(ctx)

	fa, err := fieldArgs(input.fields)
	if err != nil {
		write(errorPayload("fields_get", fmt.Errorf("invalid fields: %w", err)))
		os.Exit(1)
	}

	client, err := conn.Connect()
	if err != nil {
		write(errorPayload("fields_get", fmt.Errorf("cannot connect to Odoo: %w", err)))
		os.Exit(1)
	}

	result, err := client.ExecuteKW(input.model, "fields_get",
		godoorpc.Args{fa},
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
