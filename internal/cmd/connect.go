package cmd

import (
	"errors"
	"flag"

	"github.com/lxkrmr/godoorpc"
)

// ConnFlags holds the connection data parsed from global flags.
type ConnFlags struct {
	URL      string
	DB       string
	User     string
	Password string
}

// RegisterConnFlags registers the connection flags on a FlagSet.
func RegisterConnFlags(fs *flag.FlagSet, c *ConnFlags) {
	fs.StringVar(&c.URL, "url", "", "Odoo base URL (e.g. http://localhost:8069)")
	fs.StringVar(&c.DB, "db", "", "Database name")
	fs.StringVar(&c.User, "user", "", "Login user")
	fs.StringVar(&c.Password, "password", "", "Login password")
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
