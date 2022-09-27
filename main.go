package main

import (
	"net/http"
	"os"

	netatmo "github.com/exzz/netatmo-api-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/xperimental/netatmo-exporter/internal/collector"
	"github.com/xperimental/netatmo-exporter/internal/config"
	"github.com/xperimental/netatmo-exporter/internal/web"
)

var log = &logrus.Logger{
	Out: os.Stderr,
	Formatter: &logrus.TextFormatter{
		DisableTimestamp: true,
	},
	Level:        logrus.InfoLevel,
	ExitFunc:     os.Exit,
	ReportCaller: false,
}

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
