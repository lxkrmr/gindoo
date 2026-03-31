package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
)

// hoistFlags reorders args so that all --flag [value] pairs come before
// positional arguments. This allows flags to appear anywhere in the
// argument list without confusing Go's flag package, which stops parsing
// at the first non-flag argument.
//
// Assumption: all flags in gindoo take exactly one value (no boolean flags).
// A flag followed by another flag (or at end of args) is treated as
// a lone flag and passed through — the flag package will report the error.
func hoistFlags(args []string) []string {
	var flags, positionals []string
	i := 0
	for i < len(args) {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			if strings.Contains(arg, "=") {
				// --flag=value form: self-contained
				flags = append(flags, arg)
				i++
			} else if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				// --flag value form: consume both
				flags = append(flags, arg, args[i+1])
				i += 2
			} else {
				// --flag alone (e.g. --help)
				flags = append(flags, arg)
				i++
			}
		} else {
			positionals = append(positionals, arg)
			i++
		}
	}
	return append(flags, positionals...)
}

// parseFieldList parses an Odoo-style field list string into a []string.
// Accepts Python-style single-quoted list syntax: "['name', 'email']"
func parseFieldList(s string) ([]string, error) {
	s = strings.ReplaceAll(s, "'", "\"")
	var fields []string
	if err := json.Unmarshal([]byte(s), &fields); err != nil {
		return nil, fmt.Errorf("invalid fields %q: expected Odoo list syntax, e.g. \"['name', 'email']\"", s)
	}
	if len(fields) == 0 {
		return nil, fmt.Errorf("fields list must not be empty — use \"['name', 'email']\" syntax")
	}
	return fields, nil
}
