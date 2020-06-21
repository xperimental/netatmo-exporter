package main

import (
	"errors"
	"fmt"
	"time"

	netatmo "github.com/exzz/netatmo-api-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

const (
	envVarListenAddress       = "NETATMO_EXPORTER_ADDR"
	envVarLogLevel            = "NETATMO_LOG_LEVEL"
	envVarStaleDuration       = "NETATMO_AGE_STALE"
	envVarNetatmoClientID     = "NETATMO_CLIENT_ID"
	envVarNetatmoClientSecret = "NETATMO_CLIENT_SECRET"
	envVarNetatmoUsername     = "NETATMO_CLIENT_USERNAME"
	envVarNetatmoPassword     = "NETATMO_CLIENT_PASSWORD"

	flagListenAddress       = "addr"
	flagLogLevel            = "log-level"
	flagStaleDuration       = "age-stale"
	flagNetatmoClientID     = "client-id"
	flagNetatmoClientSecret = "client-secret"
	flagNetatmoUsername     = "username"
	flagNetatmoPassword     = "password"

	defaultStaleDuration = 30 * time.Minute
)

var (
	defaultConfig = config{
		Addr:          ":9210",
		LogLevel:      logLevel(logrus.InfoLevel),
		StaleDuration: defaultStaleDuration,
	}

	errNoBinaryName          = errors.New("need the binary name as first argument")
	errNoListenAddress       = errors.New("no listen address")
	errNoNetatmoClientID     = errors.New("need a NetAtmo client ID")
	errNoNetatmoClientSecret = errors.New("need a NetAtmo client secret")
	errNoNetatmoUsername     = errors.New("username can not be blank")
	errNoNetatmoPassword     = errors.New("password can not be blank")
)

type logLevel logrus.Level

func (l *logLevel) Type() string {
	return "level"
}

func (l *logLevel) String() string {
	return fmt.Sprintf("%s", logrus.Level(*l))
}

func (l *logLevel) Set(value string) error {
	level, err := logrus.ParseLevel(value)
	if err != nil {
		return err
	}
	*l = logLevel(level)

	return nil
}

type config struct {
	Addr          string
	LogLevel      logLevel
	StaleDuration time.Duration
	Netatmo       netatmo.Config
}

func parseConfig(args []string, getenv func(string) string) (config, error) {
	cfg := defaultConfig

	if len(args) < 1 {
		return cfg, errNoBinaryName
	}

	flagSet := pflag.NewFlagSet(args[0], pflag.ExitOnError)
	flagSet.StringVarP(&cfg.Addr, flagListenAddress, "a", cfg.Addr, "Address to listen on.")
	flagSet.Var(&cfg.LogLevel, flagLogLevel, "Sets the minimum level output through logging.")
	flagSet.DurationVar(&cfg.StaleDuration, flagStaleDuration, cfg.StaleDuration, "Data age to consider as stale. Stale data does not create metrics anymore.")
	flagSet.StringVarP(&cfg.Netatmo.ClientID, flagNetatmoClientID, "i", cfg.Netatmo.ClientID, "Client ID for NetAtmo app.")
	flagSet.StringVarP(&cfg.Netatmo.ClientSecret, flagNetatmoClientSecret, "s", cfg.Netatmo.ClientSecret, "Client secret for NetAtmo app.")
	flagSet.StringVarP(&cfg.Netatmo.Username, flagNetatmoUsername, "u", cfg.Netatmo.Username, "Username of NetAtmo account.")
	flagSet.StringVarP(&cfg.Netatmo.Password, flagNetatmoPassword, "p", cfg.Netatmo.Password, "Password of NetAtmo account.")
	flagSet.Parse(args[1:])

	if err := applyEnvironment(&cfg, getenv); err != nil {
		return config{}, fmt.Errorf("error in environment: %s", err)
	}

	if len(cfg.Addr) == 0 {
		return config{}, errNoListenAddress
	}

	if len(cfg.Netatmo.ClientID) == 0 {
		return config{}, errNoNetatmoClientID
	}

	if len(cfg.Netatmo.ClientSecret) == 0 {
		return config{}, errNoNetatmoClientSecret
	}

	if len(cfg.Netatmo.Username) == 0 {
		return config{}, errNoNetatmoUsername
	}

	if len(cfg.Netatmo.Password) == 0 {
		return config{}, errNoNetatmoPassword
	}

	return cfg, nil
}

func applyEnvironment(cfg *config, getenv func(string) string) error {
	if envAddr := getenv(envVarListenAddress); envAddr != "" {
		cfg.Addr = envAddr
	}

	if envLogLevel := getenv(envVarLogLevel); envLogLevel != "" {
		if err := cfg.LogLevel.Set(envLogLevel); err != nil {
			return err
		}
	}

	if envStaleDuration := getenv(envVarStaleDuration); envStaleDuration != "" {
		duration, err := time.ParseDuration(envStaleDuration)
		if err != nil {
			return err
		}

		cfg.StaleDuration = duration
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

	return nil
}
