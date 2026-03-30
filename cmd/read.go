package cmd

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/lxkrmr/godoorpc"
)

const readHelp = `Read fields from a single Odoo record by ID.

Usage:
  gindoo read [flags] <model> <id> <fields...>

Arguments:
  model     Technical model name (e.g. res.partner)
  id        Record ID
  fields    One or more field names to read

Connection flags:
  --url       Odoo base URL (e.g. http://localhost:8069)
  --db        Database name
  --user      Login user
  --password  Login password

Examples:
  gindoo read res.partner 1 name email phone
  gindoo read sale.order 42 name state amount_total`

// readInput holds the parsed data for a read command.
type readInput struct {
	conn   connFlags
	model  string
	id     int
	fields []string
}

// parseReadArgs parses flags and positional args — calculation.
func parseReadArgs(args []string) (readInput, error) {
	fs := flag.NewFlagSet("read", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	fs.Usage = func() { fmt.Println(readHelp) }

	var input readInput
	registerConnFlags(fs, &input.conn)

	if err := fs.Parse(args); err != nil {
		return readInput{}, err
	}

	positional := fs.Args()
	if len(positional) < 3 {
		return readInput{}, fmt.Errorf("model, id, and at least one field are required — run 'gindoo read --help'")
	}

	input.model = positional[0]

	id, err := strconv.Atoi(positional[1])
	if err != nil {
		return readInput{}, fmt.Errorf("id must be an integer, got %q", positional[1])
	}
	input.id = id
	input.fields = positional[2:]

	return input, nil
}

// buildReadResult shapes the data for the JSON response — pure calculation.
func buildReadResult(input readInput, record any) map[string]any {
	return map[string]any{
		"model":  input.model,
		"id":     input.id,
		"fields": input.fields,
		"record": record,
	}
}

// RunRead orchestrates side effects: parse, connect, execute, write.
func RunRead(args []string) {
	input, err := parseReadArgs(args)
	if err == flag.ErrHelp {
		os.Exit(0)
	}
	if err != nil {
		write(errorPayload("read", err))
		os.Exit(1)
	}

	client, err := input.conn.connect()
	if err != nil {
		write(errorPayload("read", fmt.Errorf("cannot connect to Odoo: %w", err)))
		os.Exit(1)
	}

	record, err := client.ExecuteKW(input.model, "read",
		godoorpc.Args{[]int{input.id}},
		godoorpc.KWArgs{"fields": input.fields},
	)
	if err != nil {
		write(errorPayload("read", fmt.Errorf("read failed for %s/%d: %w", input.model, input.id, err)))
		os.Exit(1)
	}

	write(successPayload("read", buildReadResult(input, record)))
}
