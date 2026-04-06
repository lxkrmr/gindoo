package cmd

import (
	"os"
	"testing"
)

func TestContextData(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	originalConfigDir := os.Getenv("HOME")
	defer os.Setenv("HOME", originalConfigDir)

	// Override config path for testing
	tmpHome := tmpDir
	os.Setenv("HOME", tmpHome)

	// Test: save and load context data
	cd := contextData{
		Contexts: map[string]connectionConfig{
			"test": {
				URL:      "http://localhost:8069",
				DB:       "testdb",
				User:     "admin",
				Password: "secret",
			},
		},
		CurrentContext: "test",
	}

	if err := saveContextData(cd); err != nil {
		t.Fatalf("could not save context data: %v", err)
	}

	loaded, err := loadContextData()
	if err != nil {
		t.Fatalf("could not load context data: %v", err)
	}

	if loaded.CurrentContext != cd.CurrentContext {
		t.Errorf("expected current context %q, got %q", cd.CurrentContext, loaded.CurrentContext)
	}

	ctx, ok := loaded.Contexts["test"]
	if !ok {
		t.Fatalf("context 'test' not found")
	}

	if ctx.URL != cd.Contexts["test"].URL {
		t.Errorf("expected URL %q, got %q", cd.Contexts["test"].URL, ctx.URL)
	}
}

func TestGetCurrentContext(t *testing.T) {
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tmpDir)

	// No context set yet
	_, _, err := GetCurrentContext()
	if err == nil {
		t.Fatal("expected error when no context is set")
	}

	// Create and set a context
	cd := contextData{
		Contexts: map[string]connectionConfig{
			"dev": {
				URL:      "http://localhost:8069",
				DB:       "dev",
				User:     "admin",
				Password: "test",
			},
		},
		CurrentContext: "dev",
	}

	if err := saveContextData(cd); err != nil {
		t.Fatalf("could not save context: %v", err)
	}

	name, ctx, err := GetCurrentContext()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if name != "dev" {
		t.Errorf("expected context name 'dev', got %q", name)
	}

	if ctx.URL != "http://localhost:8069" {
		t.Errorf("expected URL %q, got %q", "http://localhost:8069", ctx.URL)
	}
}

func TestSetAndRemoveContext(t *testing.T) {
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tmpDir)

	// Create two contexts
	cd := contextData{
		Contexts: map[string]connectionConfig{
			"dev": {
				URL:  "http://localhost:8069",
				DB:   "dev",
				User: "admin",
			},
			"staging": {
				URL:  "http://staging.example.com:8069",
				DB:   "staging",
				User: "admin",
			},
		},
		CurrentContext: "dev",
	}

	if err := saveContextData(cd); err != nil {
		t.Fatalf("could not save context: %v", err)
	}

	// Switch to staging
	if err := SetCurrentContext("staging"); err != nil {
		t.Fatalf("could not set current context: %v", err)
	}

	name, _, err := GetCurrentContext()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if name != "staging" {
		t.Errorf("expected current context 'staging', got %q", name)
	}

	// Remove staging
	if err := RemoveContext("staging"); err != nil {
		t.Fatalf("could not remove context: %v", err)
	}

	// Try to get staging, should fail
	_, err = getContext("staging")
	if err == nil {
		t.Fatal("expected error when accessing removed context")
	}

	// Current context should be cleared
	_, _, err = GetCurrentContext()
	if err == nil {
		t.Fatal("expected error when current context is removed")
	}
}

func TestConvertContextToConnFlags(t *testing.T) {
	ctx := connectionConfig{
		URL:      "http://localhost:8069",
		DB:       "mydb",
		User:     "admin",
		Password: "secret",
	}

	flags := ConvertContextToConnFlags(ctx)

	if flags.URL != ctx.URL {
		t.Errorf("expected URL %q, got %q", ctx.URL, flags.URL)
	}
	if flags.DB != ctx.DB {
		t.Errorf("expected DB %q, got %q", ctx.DB, flags.DB)
	}
	if flags.User != ctx.User {
		t.Errorf("expected User %q, got %q", ctx.User, flags.User)
	}
	if flags.Password != ctx.Password {
		t.Errorf("expected Password %q, got %q", ctx.Password, flags.Password)
	}
}
