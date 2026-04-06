package cmd

import (
	"errors"

	"github.com/lxkrmr/godoorpc"
)

// ConnFlags holds the connection data parsed from global flags.
type ConnFlags struct {
	URL      string
	DB       string
	User     string
	Password string
}

// validate checks that all connection fields are present — pure calculation.
func (c ConnFlags) validate() error {
	if c.URL == "" {
		return errors.New("--url is required (e.g. --url http://localhost:8069)")
	}
	if c.DB == "" {
		return errors.New("--db is required")
	}
	if c.User == "" {
		return errors.New("--user is required")
	}
	if c.Password == "" {
		return errors.New("--password is required")
	}
	return nil
}

// Connect validates and opens an Odoo session — side effect.
func (c ConnFlags) Connect() (*godoorpc.Client, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}
	return godoorpc.NewSession(c.URL, c.DB, c.User, c.Password)
}
