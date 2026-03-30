package cmd

import (
	"encoding/json"
	"fmt"
	"os"
)

// successPayload builds a success response — pure calculation.
func successPayload(command string, data any) map[string]any {
	return map[string]any{
		"ok":      true,
		"command": command,
		"data":    data,
	}
}

// errorPayload builds an error response — pure calculation.
func errorPayload(command string, err error) map[string]any {
	return map[string]any{
		"ok":      false,
		"command": command,
		"error":   err.Error(),
	}
}

// WriteError writes an error payload to stdout — exported for use in main.
func WriteError(command string, err error) {
	write(errorPayload(command, err))
}

// write is the single side effect: marshal and print to stdout.
func write(v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, `{"ok":false,"command":"","error":"failed to marshal output: %s"}`, err)
		os.Exit(1)
	}
	fmt.Println(string(b))
}
