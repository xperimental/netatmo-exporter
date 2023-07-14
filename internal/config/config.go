package config

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/exzz/netatmo-api-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

const (
	envVarListenAddress       = "NETATMO_EXPORTER_ADDR"
	envVarExternalURL         = "NETATMO_EXPORTER_EXTERNAL_URL"
	envVarTokenFile           = "NETATMO_EXPORTER_TOKEN_FILE"
	envVarDebugHandlers       = "DEBUG_HANDLERS"
	envVarLogLevel            = "NETATMO_LOG_LEVEL"
	envVarRefreshInterval     = "NETATMO_REFRESH_INTERVAL"
	envVarStaleDuration       = "NETATMO_AGE_STALE"
	envVarNetatmoClientID     = "NETATMO_CLIENT_ID"
	envVarNetatmoClientSecret = "NETATMO_CLIENT_SECRET"

	flagListenAddress       = "addr"
	flagExternalURL         = "external-url"
	flagTokenFile           = "token-file"
	flagDebugHandlers       = "debug-handlers"
	flagLogLevel            = "log-level"
	flagRefreshInterval     = "refresh-interval"
	flagStaleDuration       = "age-stale"
	flagNetatmoClientID     = "client-id"
	flagNetatmoClientSecret = "client-secret"

	defaultRefreshInterval = 8 * time.Minute
	defaultStaleDuration   = 60 * time.Minute
)

var (
	defaultConfig = Config{
		Addr:            ":9210",
		LogLevel:        logLevel(logrus.InfoLevel),
		RefreshInterval: defaultRefreshInterval,
		StaleDuration:   defaultStaleDuration,
	}

	errNoBinaryName          = errors.New("need the binary name as first argument")
	errNoListenAddress       = errors.New("no listen address")
	errNoNetatmoClientID     = errors.New("need a NetAtmo client ID")
	errNoNetatmoClientSecret = errors.New("need a NetAtmo client secret")
)

type logLevel logrus.Level

func (l *logLevel) Type() string {
	return "level"
}

func (l *logLevel) String() string {
	return logrus.Level(*l).String()
}

func (l *logLevel) Set(value string) error {
	level, err := logrus.ParseLevel(value)
	if err != nil {
		return err
	}
	*l = logLevel(level)

	return nil
}

// Config contains the configuration options.
type Config struct {
	Addr            string
	ExternalURL     string
	TokenFile       string
	DebugHandlers   bool
	LogLevel        logLevel
	RefreshInterval time.Duration
	StaleDuration   time.Duration
	Netatmo         netatmo.Config
}

// Parse takes the arguments and environment variables provided and creates the Config from that.
func Parse(args []string, getEnv func(string) string) (Config, error) {
	cfg := defaultConfig

	if len(args) < 1 {
		return cfg, errNoBinaryName
	}

	flagSet := pflag.NewFlagSet(args[0], pflag.ContinueOnError)
	flagSet.StringVarP(&cfg.Addr, flagListenAddress, "a", cfg.Addr, "Address to listen on.")
	flagSet.StringVar(&cfg.ExternalURL, flagExternalURL, cfg.ExternalURL, "External URL to use as base for OAuth redirect URL.")
	flagSet.StringVar(&cfg.TokenFile, flagTokenFile, cfg.TokenFile, "Path to token file for loading/persisting authentication token.")
	flagSet.BoolVar(&cfg.DebugHandlers, flagDebugHandlers, cfg.DebugHandlers, "Enables debugging HTTP handlers.")
	flagSet.Var(&cfg.LogLevel, flagLogLevel, "Sets the minimum level output through logging.")
	flagSet.DurationVar(&cfg.RefreshInterval, flagRefreshInterval, cfg.RefreshInterval, "Time interval used for internal caching of NetAtmo sensor data.")
	flagSet.DurationVar(&cfg.StaleDuration, flagStaleDuration, cfg.StaleDuration, "Data age to consider as stale. Stale data does not create metrics anymore.")
	flagSet.StringVarP(&cfg.Netatmo.ClientID, flagNetatmoClientID, "i", cfg.Netatmo.ClientID, "Client ID for NetAtmo app.")
	flagSet.StringVarP(&cfg.Netatmo.ClientSecret, flagNetatmoClientSecret, "s", cfg.Netatmo.ClientSecret, "Client secret for NetAtmo app.")

	if err := flagSet.Parse(args[1:]); err != nil {
		return Config{}, err
	}

	if err := applyEnvironment(&cfg, getEnv); err != nil {
		return Config{}, fmt.Errorf("error in environment: %s", err)
	}

	if len(cfg.Addr) == 0 {
		return Config{}, errNoListenAddress
	}

	if cfg.ExternalURL == "" {
		host, port, err := net.SplitHostPort(cfg.Addr)
		if err != nil {
			return Config{}, fmt.Errorf("error generating external URL from listen address: %w", err)
		}

		if host == "" {
			host = "127.0.0.1"
		}

		cfg.ExternalURL = fmt.Sprintf("http://%s:%s", host, port)
	}

	if len(cfg.Netatmo.ClientID) == 0 {
		return Config{}, errNoNetatmoClientID
	}

	if len(cfg.Netatmo.ClientSecret) == 0 {
		return Config{}, errNoNetatmoClientSecret
	}

	if cfg.StaleDuration < cfg.RefreshInterval {
		return Config{}, fmt.Errorf("stale duration smaller than refresh interval: %s < %s", cfg.StaleDuration, cfg.RefreshInterval)
	}

	return cfg, nil
}

func applyEnvironment(cfg *Config, getenv func(string) string) error {
	if envAddr := getenv(envVarListenAddress); envAddr != "" {
		cfg.Addr = envAddr
	}

	if externalURL := getenv(envVarExternalURL); externalURL != "" {
		cfg.ExternalURL = externalURL
	}

	if tokenFile := getenv(envVarTokenFile); tokenFile != "" {
		cfg.TokenFile = tokenFile
	}

	if envDebugHandlers := getenv(envVarDebugHandlers); envDebugHandlers != "" {
		cfg.DebugHandlers = true
	}

	if envLogLevel := getenv(envVarLogLevel); envLogLevel != "" {
		if err := cfg.LogLevel.Set(envLogLevel); err != nil {
			return err
		}
	}

	if envRefreshInterval := getenv(envVarRefreshInterval); envRefreshInterval != "" {
		duration, err := time.ParseDuration(envRefreshInterval)
		if err != nil {
			return err
		}

		cfg.RefreshInterval = duration
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

	return nil
}
