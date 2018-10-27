package main

import (
	"errors"

	netatmo "github.com/exzz/netatmo-api-go"
	"github.com/spf13/pflag"
)

const (
	envVarListenAddress       = "NETATMO_EXPORTER_ADDR"
	envVarNetatmoClientID     = "NETATMO_CLIENT_ID"
	envVarNetatmoClientSecret = "NETATMO_CLIENT_SECRET"
	envVarNetatmoUsername     = "NETATMO_CLIENT_USERNAME"
	envVarNetatmoPassword     = "NETATMO_CLIENT_PASSWORD"

	flagListenAddress       = "addr"
	flagNetatmoClientID     = "client-id"
	flagNetatmoClientSecret = "client-secret"
	flagNetatmoUsername     = "username"
	flagNetatmoPassword     = "password"
)

var (
	defaultConfig = config{
		Addr: ":9210",
	}

	errNoBinaryName          = errors.New("need the binary name as first argument")
	errNoListenAddress       = errors.New("no listen address")
	errNoNetatmoClientID     = errors.New("need a NetAtmo client ID")
	errNoNetatmoClientSecret = errors.New("need a NetAtmo client secret")
	errNoNetatmoUsername     = errors.New("username can not be blank")
	errNoNetatmoPassword     = errors.New("password can not be blank")
)

type config struct {
	Addr    string
	Netatmo netatmo.Config
}

func parseConfig(args []string, getenv func(string) string) (config, error) {
	cfg := defaultConfig

	if len(args) < 1 {
		return cfg, errNoBinaryName
	}

	flagSet := pflag.NewFlagSet(args[0], pflag.ExitOnError)
	flagSet.StringVarP(&cfg.Addr, "addr", "a", cfg.Addr, "Address to listen on.")
	flagSet.StringVarP(&cfg.Netatmo.ClientID, "client-id", "i", cfg.Netatmo.ClientID, "Client ID for NetAtmo app.")
	flagSet.StringVarP(&cfg.Netatmo.ClientSecret, "client-secret", "s", cfg.Netatmo.ClientSecret, "Client secret for NetAtmo app.")
	flagSet.StringVarP(&cfg.Netatmo.Username, "username", "u", cfg.Netatmo.Username, "Username of NetAtmo account.")
	flagSet.StringVarP(&cfg.Netatmo.Password, "password", "p", cfg.Netatmo.Password, "Password of NetAtmo account.")
	flagSet.Parse(args[1:])

	applyEnvironment(&cfg, getenv)

	if len(cfg.Addr) == 0 {
		return cfg, errNoListenAddress
	}

	if len(cfg.Netatmo.ClientID) == 0 {
		return cfg, errNoNetatmoClientID
	}

	if len(cfg.Netatmo.ClientSecret) == 0 {
		return cfg, errNoNetatmoClientSecret
	}

	if len(cfg.Netatmo.Username) == 0 {
		return cfg, errNoNetatmoUsername
	}

	if len(cfg.Netatmo.Password) == 0 {
		return cfg, errNoNetatmoPassword
	}

	return cfg, nil
}

func applyEnvironment(cfg *config, getenv func(string) string) {
	if envAddr := getenv(envVarListenAddress); envAddr != "" {
		cfg.Addr = envAddr
	}

	if envClientID := getenv(envVarNetatmoClientID); envClientID != "" {
		cfg.Netatmo.ClientID = envClientID
	}

	if envClientSecret := getenv(envVarNetatmoClientSecret); envClientSecret != "" {
		cfg.Netatmo.ClientSecret = envClientSecret
	}

	if envUsername := getenv(envVarNetatmoUsername); envUsername != "" {
		cfg.Netatmo.Username = envUsername
	}

	if envPassword := getenv(envVarNetatmoPassword); envPassword != "" {
		cfg.Netatmo.Password = envPassword
	}
}
