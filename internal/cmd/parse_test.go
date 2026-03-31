package cmd

import (
	"testing"
)

func TestParseSearchReadArgs(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		expectModel  string
		expectDomain string
		expectFields string
		expectLimit  int
		wantErr      bool
	}{
		{
			name:         "model domain fields",
			args:         []string{"res.partner", "[]", "['name', 'email']"},
			expectModel:  "res.partner",
			expectDomain: "[]",
			expectFields: "['name', 'email']",
			expectLimit:  10,
		},
		{
			name:         "with --limit before positionals",
			args:         []string{"--limit", "5", "res.partner", "[]", "['name']"},
			expectModel:  "res.partner",
			expectDomain: "[]",
			expectFields: "['name']",
			expectLimit:  5,
		},
		{
			name:         "with --limit after positionals",
			args:         []string{"res.partner", "[]", "['name']", "--limit", "5"},
			expectModel:  "res.partner",
			expectDomain: "[]",
			expectFields: "['name']",
			expectLimit:  5,
		},
		{
			name:         "with --limit between positionals",
			args:         []string{"res.partner", "[]", "--limit", "20", "['name']"},
			expectModel:  "res.partner",
			expectDomain: "[]",
			expectFields: "['name']",
			expectLimit:  20,
		},
		{
			name:    "missing fields",
			args:    []string{"res.partner", "[]"},
			wantErr: true,
		},
		{
			name:    "missing domain and fields",
			args:    []string{"res.partner"},
			wantErr: true,
		},
		{
			name:    "missing all",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "too many positionals",
			args:    []string{"res.partner", "[]", "['name']", "extra"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSearchReadArgs(tt.args)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.model != tt.expectModel {
				t.Errorf("model: expected %q, got %q", tt.expectModel, got.model)
			}
			if got.domain != tt.expectDomain {
				t.Errorf("domain: expected %q, got %q", tt.expectDomain, got.domain)
			}
			if got.fields != tt.expectFields {
				t.Errorf("fields: expected %q, got %q", tt.expectFields, got.fields)
			}
			if got.limit != tt.expectLimit {
				t.Errorf("limit: expected %d, got %d", tt.expectLimit, got.limit)
			}
		})
	}
}

func TestParseSearchCountArgs(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		expectModel  string
		expectDomain string
		wantErr      bool
	}{
		{
			name:         "model and domain",
			args:         []string{"res.partner", "[]"},
			expectModel:  "res.partner",
			expectDomain: "[]",
		},
		{
			name:         "with domain filter",
			args:         []string{"res.partner", "[('is_company', '=', True)]"},
			expectModel:  "res.partner",
			expectDomain: "[('is_company', '=', True)]",
		},
		{
			name:    "missing domain",
			args:    []string{"res.partner"},
			wantErr: true,
		},
		{
			name:    "missing all",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "too many positionals",
			args:    []string{"res.partner", "[]", "extra"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSearchCountArgs(tt.args)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.model != tt.expectModel {
				t.Errorf("model: expected %q, got %q", tt.expectModel, got.model)
			}
			if got.domain != tt.expectDomain {
				t.Errorf("domain: expected %q, got %q", tt.expectDomain, got.domain)
			}
		})
	}
}

func TestParseFieldsGetArgs(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		expectModel  string
		expectFields string
		wantErr      bool
	}{
		{
			name:        "model only",
			args:        []string{"res.partner"},
			expectModel: "res.partner",
		},
		{
			name:         "model with fields",
			args:         []string{"res.partner", "['name', 'email']"},
			expectModel:  "res.partner",
			expectFields: "['name', 'email']",
		},
		{
			name:    "missing model",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "too many positionals",
			args:    []string{"res.partner", "['name']", "extra"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFieldsGetArgs(tt.args)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.model != tt.expectModel {
				t.Errorf("model: expected %q, got %q", tt.expectModel, got.model)
			}
			if got.fields != tt.expectFields {
				t.Errorf("fields: expected %q, got %q", tt.expectFields, got.fields)
			}
		})
	}
}

func TestParseFieldList(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		expect  []string
		wantErr bool
	}{
		{
			name:   "single field",
			input:  "['name']",
			expect: []string{"name"},
		},
		{
			name:   "multiple fields",
			input:  "['name', 'email', 'phone']",
			expect: []string{"name", "email", "phone"},
		},
		{
			name:    "empty list",
			input:   "[]",
			wantErr: true,
		},
		{
			name:    "invalid syntax",
			input:   "name,email",
			wantErr: true,
		},
		{
			name:    "plain string",
			input:   "name",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFieldList(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != len(tt.expect) {
				t.Fatalf("expected %v, got %v", tt.expect, got)
			}
			for i, f := range tt.expect {
				if got[i] != f {
					t.Errorf("field[%d]: expected %q, got %q", i, f, got[i])
				}
			}
		})
	}
}

func TestFieldArgs(t *testing.T) {
	// empty string → false (all fields)
	got, err := fieldArgs("")
	if err != nil {
		t.Fatalf("unexpected error for empty fields: %v", err)
	}
	if got != false {
		t.Errorf("expected false for empty fields, got %v", got)
	}

	// valid list → []string
	got, err = fieldArgs("['name', 'email']")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fields, ok := got.([]string)
	if !ok || len(fields) != 2 {
		t.Errorf("expected []string with 2 elements, got %v", got)
	}

	// invalid → error
	_, err = fieldArgs("[]")
	if err == nil {
		t.Error("expected error for empty list, got nil")
	}
}

func TestHoistFlags(t *testing.T) {
	tests := []struct {
		name   string
		input  []string
		expect []string
	}{
		{
			name:   "no flags",
			input:  []string{"res.partner", "[]", "['name']"},
			expect: []string{"res.partner", "[]", "['name']"},
		},
		{
			name:   "flag before positionals",
			input:  []string{"--limit", "5", "res.partner", "[]", "['name']"},
			expect: []string{"--limit", "5", "res.partner", "[]", "['name']"},
		},
		{
			name:   "flag after positionals",
			input:  []string{"res.partner", "[]", "['name']", "--limit", "5"},
			expect: []string{"--limit", "5", "res.partner", "[]", "['name']"},
		},
		{
			name:   "flag between positionals",
			input:  []string{"res.partner", "[]", "--limit", "5", "['name']"},
			expect: []string{"--limit", "5", "res.partner", "[]", "['name']"},
		},
		{
			name:   "flag=value form",
			input:  []string{"res.partner", "[]", "['name']", "--limit=5"},
			expect: []string{"--limit=5", "res.partner", "[]", "['name']"},
		},
		{
			name:   "help flag alone",
			input:  []string{"--help"},
			expect: []string{"--help"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hoistFlags(tt.input)
			if len(got) != len(tt.expect) {
				t.Fatalf("expected %v, got %v", tt.expect, got)
			}
			for i, v := range tt.expect {
				if got[i] != v {
					t.Errorf("[%d]: expected %q, got %q", i, v, got[i])
				}
			}
		})
	}
}
