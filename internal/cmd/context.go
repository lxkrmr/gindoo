package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// contextData holds all saved contexts and the current one.
type contextData struct {
	Contexts       map[string]connectionConfig `json:"contexts"`
	CurrentContext string                      `json:"current_context"`
}

// connectionConfig holds URL, DB, User, Password for a context.
type connectionConfig struct {
	URL      string `json:"url"`
	DB       string `json:"db"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// getConfigPath returns the path to the contexts.json file.
// Uses ~/.config/gindoo/ on Unix, %APPDATA%\gindoo\ on Windows.
func getConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("could not determine config directory: %w", err)
	}

	gindooDir := filepath.Join(configDir, "gindoo")
	return filepath.Join(gindooDir, "contexts.json"), nil
}

// loadContextData loads the contexts.json file, returns empty data if it doesn't exist.
func loadContextData() (contextData, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return contextData{}, err
	}

	data, err := os.ReadFile(configPath)
	if errors.Is(err, os.ErrNotExist) {
		return contextData{Contexts: make(map[string]connectionConfig)}, nil
	}
	if err != nil {
		return contextData{}, fmt.Errorf("could not read config file %q: %w", configPath, err)
	}

	var cd contextData
	if err := json.Unmarshal(data, &cd); err != nil {
		return contextData{}, fmt.Errorf("could not parse config file: %w", err)
	}

	if cd.Contexts == nil {
		cd.Contexts = make(map[string]connectionConfig)
	}

	return cd, nil
}

// saveContextData writes contexts.json, creating directories as needed.
func saveContextData(cd contextData) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("could not create directory %q: %w", configDir, err)
	}

	data, err := json.MarshalIndent(cd, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal contexts: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("could not write config file %q: %w", configPath, err)
	}

	return nil
}

// readPassword reads a password from stdin (visible in first iteration).
func readPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	var password string
	_, err := fmt.Scanln(&password)
	if err != nil {
		return "", fmt.Errorf("could not read password: %w", err)
	}
	return password, nil
}

// readInput reads a line from stdin.
func readInput(prompt string) (string, error) {
	fmt.Print(prompt)
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return "", fmt.Errorf("could not read input: %w", err)
	}
	return input, nil
}

// getContext loads a context by name.
// Returns error if context doesn't exist.
func getContext(name string) (connectionConfig, error) {
	cd, err := loadContextData()
	if err != nil {
		return connectionConfig{}, err
	}

	ctx, ok := cd.Contexts[name]
	if !ok {
		return connectionConfig{}, fmt.Errorf("context %q not found", name)
	}

	return ctx, nil
}

// GetCurrentContext returns the current context, or an error if not set.
func GetCurrentContext() (string, connectionConfig, error) {
	cd, err := loadContextData()
	if err != nil {
		return "", connectionConfig{}, err
	}

	if cd.CurrentContext == "" {
		return "", connectionConfig{},
			errors.New("no context set - run 'gindoo context create <name>' and 'gindoo context use <name>'")
	}

	ctx, ok := cd.Contexts[cd.CurrentContext]
	if !ok {
		return "", connectionConfig{},
			fmt.Errorf("current context %q not found in contexts", cd.CurrentContext)
	}

	return cd.CurrentContext, ctx, nil
}

// CreateContextInteractive runs the context create wizard.
func CreateContextInteractive(name string) error {
	if name == "" {
		return errors.New("context name is required")
	}

	cd, err := loadContextData()
	if err != nil {
		return err
	}

	if _, exists := cd.Contexts[name]; exists {
		return fmt.Errorf("context %q already exists", name)
	}

	url, err := readInput("URL (e.g. http://localhost:8069): ")
	if err != nil {
		return err
	}
	if url == "" {
		return errors.New("URL is required")
	}

	db, err := readInput("Database: ")
	if err != nil {
		return err
	}
	if db == "" {
		return errors.New("Database is required")
	}

	user, err := readInput("User: ")
	if err != nil {
		return err
	}
	if user == "" {
		return errors.New("User is required")
	}

	password, err := readPassword("Password: ")
	if err != nil {
		return err
	}
	if password == "" {
		return errors.New("Password is required")
	}

	cd.Contexts[name] = connectionConfig{
		URL:      url,
		DB:       db,
		User:     user,
		Password: password,
	}

	// Set as current context if it's the first one
	if cd.CurrentContext == "" {
		cd.CurrentContext = name
	}

	if err := saveContextData(cd); err != nil {
		return err
	}

	return nil
}

// ListContexts returns all context names with the current one marked.
func ListContexts() ([]string, string, error) {
	cd, err := loadContextData()
	if err != nil {
		return nil, "", err
	}

	names := make([]string, 0, len(cd.Contexts))
	for name := range cd.Contexts {
		names = append(names, name)
	}

	return names, cd.CurrentContext, nil
}

// SetCurrentContext sets the current context.
func SetCurrentContext(name string) error {
	cd, err := loadContextData()
	if err != nil {
		return err
	}

	if _, exists := cd.Contexts[name]; !exists {
		return fmt.Errorf("context %q not found", name)
	}

	cd.CurrentContext = name
	return saveContextData(cd)
}

// RemoveContext deletes a context.
func RemoveContext(name string) error {
	cd, err := loadContextData()
	if err != nil {
		return err
	}

	if _, exists := cd.Contexts[name]; !exists {
		return fmt.Errorf("context %q not found", name)
	}

	delete(cd.Contexts, name)

	// If we deleted the current context, clear it
	if cd.CurrentContext == name {
		cd.CurrentContext = ""
	}

	return saveContextData(cd)
}

// ConvertContextToConnFlags converts a connectionConfig to ConnFlags for use with Connect().
func ConvertContextToConnFlags(ctx connectionConfig) ConnFlags {
	return ConnFlags{
		URL:      ctx.URL,
		DB:       ctx.DB,
		User:     ctx.User,
		Password: ctx.Password,
	}
}
