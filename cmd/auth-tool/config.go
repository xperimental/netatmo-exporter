package main

import (
	"errors"

	"github.com/spf13/pflag"
)

type Config struct {
	ClientID     string
	ClientSecret string
}

func parseArgs(cmd string, args []string) (*Config, error) {
	cfg := &Config{}

	flags := pflag.NewFlagSet(cmd, pflag.ContinueOnError)
	flags.StringVar(&cfg.ClientID, "client-id", cfg.ClientID, "Client ID of OAuth application.")
	flags.StringVar(&cfg.ClientSecret, "client-secret", cfg.ClientSecret, "Client secret of OAuth application.")
	if err := flags.Parse(args); err != nil {
		log.Fatalf("Error in flags: %s", err)
	}

	if cfg.ClientID == "" {
		return nil, errors.New("need client-id")
	}

	if cfg.ClientSecret == "" {
		return nil, errors.New("need client-secret")
	}

	return cfg, nil
}
