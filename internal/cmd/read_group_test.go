package cmd

import (
	"testing"
)

func TestParseReadGroupArgs(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		expectModel  string
		expectDomain string
		expectFields string
		expectGroupBy string
		expectLimit  int
		expectOrderBy string
		wantErr      bool
	}{
		{
			name:         "model domain fields groupby",
			args:         []string{"product.template", "[]", "['fine_weight:avg']", "['default_code']"},
			expectModel:  "product.template",
			expectDomain: "[]",
			expectFields: "['fine_weight:avg']",
			expectGroupBy: "['default_code']",
			expectLimit:  10,
		},
		{
			name:         "with --limit before positionals",
			args:         []string{"--limit", "5", "product.template", "[]", "['fine_weight']", "['default_code']"},
			expectModel:  "product.template",
			expectDomain: "[]",
			expectFields: "['fine_weight']",
			expectGroupBy: "['default_code']",
			expectLimit:  5,
		},
		{
			name:         "with --limit after positionals",
			args:         []string{"product.template", "[]", "['fine_weight']", "['default_code']", "--limit", "5"},
			expectModel:  "product.template",
			expectDomain: "[]",
			expectFields: "['fine_weight']",
			expectGroupBy: "['default_code']",
			expectLimit:  5,
		},
		{
			name:         "with --orderby",
			args:         []string{"product.template", "[]", "['fine_weight']", "['default_code']", "--orderby", "default_code desc"},
			expectModel:  "product.template",
			expectDomain: "[]",
			expectFields: "['fine_weight']",
			expectGroupBy: "['default_code']",
			expectOrderBy: "default_code desc",
			expectLimit:  10,
		},
		{
			name:         "with both --limit and --orderby",
			args:         []string{"product.template", "[]", "['fine_weight']", "['default_code']", "--limit", "20", "--orderby", "default_code asc"},
			expectModel:  "product.template",
			expectDomain: "[]",
			expectFields: "['fine_weight']",
			expectGroupBy: "['default_code']",
			expectLimit:  20,
			expectOrderBy: "default_code asc",
		},
		{
			name:    "missing groupby",
			args:    []string{"product.template", "[]", "['fine_weight']"},
			wantErr: true,
		},
		{
			name:    "missing fields and groupby",
			args:    []string{"product.template", "[]"},
			wantErr: true,
		},
		{
			name:    "missing domain, fields and groupby",
			args:    []string{"product.template"},
			wantErr: true,
		},
		{
			name:    "missing all",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "too many positionals",
			args:    []string{"product.template", "[]", "['fine_weight']", "['default_code']", "extra"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseReadGroupArgs(tt.args)
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
			if got.groupby != tt.expectGroupBy {
				t.Errorf("groupby: expected %q, got %q", tt.expectGroupBy, got.groupby)
			}
			if got.limit != tt.expectLimit {
				t.Errorf("limit: expected %d, got %d", tt.expectLimit, got.limit)
			}
			if got.orderby != tt.expectOrderBy {
				t.Errorf("orderby: expected %q, got %q", tt.expectOrderBy, got.orderby)
			}
		})
	}
}

func TestBuildReadGroupResult(t *testing.T) {
	input := readGroupInput{
		model:   "product.template",
		domain:  "[]",
		fields:  "['fine_weight:avg']",
		groupby: "['default_code']",
		limit:   10,
		orderby: "default_code desc",
	}

	records := []map[string]any{
		{
			"default_code": "1001v",
			"fine_weight":  15.55174,
			"__count":      2,
		},
	}

	result := buildReadGroupResult(input, records)

	if result["model"] != "product.template" {
		t.Errorf("model: expected %q, got %q", "product.template", result["model"])
	}
	if result["domain"] != "[]" {
		t.Errorf("domain: expected %q, got %q", "[]", result["domain"])
	}
	if result["fields"] != "['fine_weight:avg']" {
		t.Errorf("fields: expected %q, got %q", "['fine_weight:avg']", result["fields"])
	}
	if result["groupby"] != "['default_code']" {
		t.Errorf("groupby: expected %q, got %q", "['default_code']", result["groupby"])
	}
	if result["limit"] != 10 {
		t.Errorf("limit: expected %d, got %d", 10, result["limit"])
	}
	if result["orderby"] != "default_code desc" {
		t.Errorf("orderby: expected %q, got %q", "default_code desc", result["orderby"])
	}
	if result["records"] == nil {
		t.Error("records should not be nil")
	}
}