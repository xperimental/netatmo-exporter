package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	netatmo "github.com/exzz/netatmo-api-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/xperimental/netatmo-exporter/internal/collector"
	"github.com/xperimental/netatmo-exporter/internal/config"
	"github.com/xperimental/netatmo-exporter/internal/logger"
	"github.com/xperimental/netatmo-exporter/internal/web"
	"golang.org/x/oauth2"
)

var (
	signals = []os.Signal{
		syscall.SIGINT,
		syscall.SIGTERM,
	}

	log = logger.NewLogger()
)

func main() {
	cfg, err := config.Parse(os.Args, os.Getenv)
	switch {
	case err == pflag.ErrHelp:
		return
	case err != nil:
		log.Fatalf("Error in configuration: %s", err)
	default:
	}
	log.SetLevel(logrus.Level(cfg.LogLevel))

	client := netatmo.NewClient(cfg.Netatmo)

	if cfg.TokenFile != "" {
		token, err := loadToken(cfg.TokenFile)
		switch {
		case os.IsNotExist(err):
		case err != nil:
			log.Fatalf("Error loading token: %s", err)
		default:
			client.InitWithToken(context.Background(), token)
		}

		registerSignalHandler(client, cfg.TokenFile)
	}

	metrics := &collector.NetatmoCollector{
		Log:             log,
		ReadFunction:    client.Read,
		RefreshInterval: cfg.RefreshInterval,
		StaleThreshold:  cfg.StaleDuration,
	}
	prometheus.MustRegister(metrics)

	if cfg.DebugHandlers {
		http.Handle("/debug/data", web.DebugHandler(log, client.Read))
	}

	http.Handle("/authorize", web.AuthorizeHandler(cfg.ExternalURL, client))
	http.Handle("/callback", web.CallbackHandler(client))
	http.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))
	http.Handle("/version", versionHandler(log))
	http.Handle("/", http.RedirectHandler("/metrics", http.StatusFound))

	log.Infof("Listen on %s...", cfg.Addr)
	log.Fatal(http.ListenAndServe(cfg.Addr, nil))
}

func loadToken(fileName string) (*oauth2.Token, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var token oauth2.Token
	if err := json.NewDecoder(file).Decode(&token); err != nil {
		return nil, err
	}

	return &token, nil
}

func registerSignalHandler(client *netatmo.Client, fileName string) {
	ch := make(chan os.Signal)
	signal.Notify(ch, signals...)
	go func() {
		sig := <-ch
		signal.Reset(signals...)
		log.Debugf("Got signal: %s", sig)

		if err := saveToken(client, fileName); err != nil {
			log.Errorf("Error persisting token: %s", err)
		}

		os.Exit(0)
	}()
}

func saveToken(client *netatmo.Client, fileName string) error {
	token, err := client.CurrentToken()
	if err != nil {
		return fmt.Errorf("error retrieving token: %w", err)
	}

	data, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("error marshalling token: %w", err)
	}

	if err := os.WriteFile(fileName, data, 0600); err != nil {
		return fmt.Errorf("error writing token file: %w", err)
	}

	return nil
}
