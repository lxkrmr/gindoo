package cmd

import (
	"testing"
)

func TestParseSearchArgs(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectModel string
		expectFields []string
		expectLimit  int
		expectOffset int
		expectDomain string
		wantErr      bool
	}{
		{
			name:         "model only",
			args:         []string{"res.partner"},
			expectModel:  "res.partner",
			expectFields: []string{"id"},
			expectLimit:  10,
			expectOffset: 0,
		},
		{
			name:         "model with fields",
			args:         []string{"res.partner", "name", "email"},
			expectModel:  "res.partner",
			expectFields: []string{"name", "email"},
		},
		{
			name:         "with domain and limit",
			args:         []string{"--domain", "[('is_company', '=', True)]", "--limit", "5", "res.partner", "name"},
			expectModel:  "res.partner",
			expectFields: []string{"name"},
			expectDomain: "[('is_company', '=', True)]",
			expectLimit:  5,
		},
		{
			name:    "missing model",
			args:    []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSearchArgs(tt.args)
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
			if tt.expectDomain != "" && got.domain != tt.expectDomain {
				t.Errorf("domain: expected %q, got %q", tt.expectDomain, got.domain)
			}
			if tt.expectLimit != 0 && got.limit != tt.expectLimit {
				t.Errorf("limit: expected %d, got %d", tt.expectLimit, got.limit)
			}
		})
	}
}

func TestParseReadArgs(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectModel string
		expectID    int
		expectFields []string
		wantErr     bool
	}{
		{
			name:         "valid",
			args:         []string{"res.partner", "1", "name", "email"},
			expectModel:  "res.partner",
			expectID:     1,
			expectFields: []string{"name", "email"},
		},
		{
			name:    "missing fields",
			args:    []string{"res.partner", "1"},
			wantErr: true,
		},
		{
			name:    "non-integer id",
			args:    []string{"res.partner", "abc", "name"},
			wantErr: true,
		},
		{
			name:    "missing model",
			args:    []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseReadArgs(tt.args)
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
			if got.id != tt.expectID {
				t.Errorf("id: expected %d, got %d", tt.expectID, got.id)
			}
		})
	}
}

func TestParseFieldsGetArgs(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectModel string
		expectFields []string
		wantErr     bool
	}{
		{
			name:        "model only",
			args:        []string{"res.partner"},
			expectModel: "res.partner",
		},
		{
			name:         "model with specific fields",
			args:         []string{"res.partner", "name", "email"},
			expectModel:  "res.partner",
			expectFields: []string{"name", "email"},
		},
		{
			name:    "missing model",
			args:    []string{},
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
			if len(tt.expectFields) > 0 {
				if len(got.fields) != len(tt.expectFields) {
					t.Errorf("fields: expected %v, got %v", tt.expectFields, got.fields)
				}
			}
		})
	}
}

func TestFieldArgs(t *testing.T) {
	if fieldArgs(nil) != false {
		t.Error("expected false for nil fields")
	}
	if fieldArgs([]string{}) != false {
		t.Error("expected false for empty fields")
	}
	fields := []string{"name", "email"}
	got, ok := fieldArgs(fields).([]string)
	if !ok || len(got) != 2 {
		t.Errorf("expected []string with 2 elements, got %v", fieldArgs(fields))
	}
}
