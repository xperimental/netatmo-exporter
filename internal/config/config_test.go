package config

import (
	"reflect"
	"testing"
	"time"

	netatmo "github.com/exzz/netatmo-api-go"
	"github.com/sirupsen/logrus"
)

func TestParseConfig(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		env        map[string]string
		wantConfig Config
		wantErr    error
	}{
		{
			name:       "no args",
			args:       []string{},
			env:        map[string]string{},
			wantConfig: Config{},
			wantErr:    errNoBinaryName,
		},
		{
			name: "success",
			args: []string{
				"test-cmd",
				"--" + flagNetatmoClientID,
				"id",
				"--" + flagNetatmoClientSecret,
				"secret",
			},
			env: map[string]string{},
			wantConfig: Config{
				Addr:            defaultConfig.Addr,
				ExternalURL:     "http://127.0.0.1:9210",
				LogLevel:        logLevel(logrus.InfoLevel),
				RefreshInterval: defaultRefreshInterval,
				StaleDuration:   defaultStaleDuration,
				Netatmo: netatmo.Config{
					ClientID:     "id",
					ClientSecret: "secret",
				},
			},
			wantErr: nil,
		},
		{
			name: "all env",
			args: []string{
				"test-cmd",
			},
			env: map[string]string{
				envVarListenAddress:       ":8080",
				envVarExternalURL:         "http://example.com",
				envVarTokenFile:           "token.json",
				envVarLogLevel:            "debug",
				envVarRefreshInterval:     "5m",
				envVarStaleDuration:       "10m",
				envVarNetatmoClientID:     "id",
				envVarNetatmoClientSecret: "secret",
			},
			wantConfig: Config{
				Addr:            ":8080",
				ExternalURL:     "http://example.com",
				TokenFile:       "token.json",
				LogLevel:        logLevel(logrus.DebugLevel),
				RefreshInterval: 5 * time.Minute,
				StaleDuration:   10 * time.Minute,
				Netatmo: netatmo.Config{
					ClientID:     "id",
					ClientSecret: "secret",
				},
			},
			wantErr: nil,
		},
		{
			name: "no addr",
			args: []string{
				"test-cmd",
				"--" + flagListenAddress,
				"",
				"--" + flagNetatmoClientID,
				"id",
				"--" + flagNetatmoClientSecret,
				"secret",
			},
			env: map[string]string{},
			wantConfig: Config{
				Addr: defaultConfig.Addr,
				Netatmo: netatmo.Config{
					ClientID:     "id",
					ClientSecret: "secret",
				},
			},
			wantErr: errNoListenAddress,
		},
		{
			name: "no client id",
			args: []string{
				"test-cmd",
				"--" + flagNetatmoClientSecret,
				"secret",
			},
			env: map[string]string{},
			wantConfig: Config{
				Addr: defaultConfig.Addr,
				Netatmo: netatmo.Config{
					ClientID:     "id",
					ClientSecret: "secret",
				},
			},
			wantErr: errNoNetatmoClientID,
		},
		{
			name: "no client secret",
			args: []string{
				"test-cmd",
				"--" + flagNetatmoClientID,
				"id",
			},
			env: map[string]string{},
			wantConfig: Config{
				Addr: defaultConfig.Addr,
				Netatmo: netatmo.Config{
					ClientID:     "id",
					ClientSecret: "secret",
				},
			},
			wantErr: errNoNetatmoClientSecret,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			getenv := func(key string) string {
				return tt.env[key]
			}

			config, err := Parse(tt.args, getenv)

			if err != tt.wantErr {
				t.Errorf("got error %q, want %q", err, tt.wantErr)
			}

			if err != nil {
				return
			}

			if !reflect.DeepEqual(config, tt.wantConfig) {
				t.Errorf("got config %v, want %v", config, tt.wantConfig)
			}
		})
	}
}
