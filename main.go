package main

import (
	"net/http"
	"os"
	"time"

	netatmo "github.com/exzz/netatmo-api-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/xperimental/netatmo-exporter/internal/collector"
	"github.com/xperimental/netatmo-exporter/internal/config"
)

var (
	log = &logrus.Logger{
		Out: os.Stderr,
		Formatter: &logrus.TextFormatter{
			DisableTimestamp: true,
		},
		Level:        logrus.InfoLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}
)

func main() {
	cfg, err := config.Parse(os.Args, os.Getenv)
	if err != nil {
		log.Fatalf("Error in configuration: %s", err)
	}
	log.SetLevel(logrus.Level(cfg.LogLevel))

	log.Infof("Login as %s", cfg.Netatmo.Username)
	client, err := netatmo.NewClient(cfg.Netatmo)
	if err != nil {
		log.Fatalf("Error creating client: %s", err)
	}

	metrics := &collector.NetatmoCollector{
		Log:             log,
		ReadFunction:    client.Read,
		RefreshInterval: cfg.RefreshInterval,
		StaleThreshold:  cfg.StaleDuration,
	}
	prometheus.MustRegister(metrics)

	// Trigger first refresh
	metrics.RefreshData(time.Now())

	http.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))
	http.Handle("/version", versionHandler(log))
	http.Handle("/", http.RedirectHandler("/metrics", http.StatusFound))

	log.Infof("Listen on %s...", cfg.Addr)
	log.Fatal(http.ListenAndServe(cfg.Addr, nil))
}
