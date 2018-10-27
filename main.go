package main

import (
	"log"
	"net/http"
	"os"

	netatmo "github.com/exzz/netatmo-api-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg, err := parseConfig(os.Args, os.Getenv)
	if err != nil {
		log.Fatalf("Error in configuration: %s", err)
	}

	log.Printf("Login as %s", cfg.Netatmo.Username)
	client, err := netatmo.NewClient(cfg.Netatmo)
	if err != nil {
		log.Fatalf("Error creating client: %s", err)
	}

	metrics := &netatmoCollector{
		client: client,
	}
	prometheus.MustRegister(metrics)

	http.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))
	http.Handle("/", http.RedirectHandler("/metrics", http.StatusFound))

	log.Printf("Listen on %s...", cfg.Addr)
	log.Fatal(http.ListenAndServe(cfg.Addr, nil))
}
