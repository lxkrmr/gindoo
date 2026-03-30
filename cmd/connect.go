package cmd

import (
	"errors"
	"flag"

	"github.com/lxkrmr/godoorpc"
)

// connFlags holds the connection data parsed from flags.
type connFlags struct {
	url      string
	db       string
	user     string
	password string
}

// registerConnFlags registers the connection flags on a FlagSet.
func registerConnFlags(fs *flag.FlagSet, c *connFlags) {
	fs.StringVar(&c.url, "url", "", "Odoo base URL (e.g. http://localhost:8069)")
	fs.StringVar(&c.db, "db", "", "Database name")
	fs.StringVar(&c.user, "user", "", "Login user")
	fs.StringVar(&c.password, "password", "", "Login password")
}

// validate checks that all connection fields are present — pure calculation.
func (c connFlags) validate() error {
	if c.url == "" {
		return errors.New("--url is required (e.g. --url http://localhost:8069)")
	}
	if c.db == "" {
		return errors.New("--db is required")
	}
	if c.user == "" {
		return errors.New("--user is required")
	}
	if c.password == "" {
		return errors.New("--password is required")
	}
	return nil
}

// connect validates and opens an Odoo session — side effect.
func (c connFlags) connect() (*godoorpc.Client, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}
	return godoorpc.NewSession(c.url, c.db, c.user, c.password)
}
