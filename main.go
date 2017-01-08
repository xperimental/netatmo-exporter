package main

import (
	"errors"
	"log"
	"net/http"

	netatmo "github.com/exzz/netatmo-api-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/pflag"
)

type config struct {
	Addr    string
	Netatmo netatmo.Config
}

func parseConfig() (config, error) {
	cfg := config{}
	pflag.StringVarP(&cfg.Addr, "addr", "a", ":8080", "Address to listen on.")
	pflag.StringVarP(&cfg.Netatmo.ClientID, "client-id", "i", "", "Client ID for NetAtmo app.")
	pflag.StringVarP(&cfg.Netatmo.ClientSecret, "client-secret", "s", "", "Client secret for NetAtmo app.")
	pflag.StringVarP(&cfg.Netatmo.Username, "username", "u", "", "Username of NetAtmo account.")
	pflag.StringVarP(&cfg.Netatmo.Password, "password", "p", "", "Password of NetAtmo account.")
	pflag.Parse()

	if len(cfg.Addr) == 0 {
		return cfg, errors.New("no listen address")
	}

	if len(cfg.Netatmo.ClientID) == 0 {
		return cfg, errors.New("need a NetAtmo client ID")
	}

	if len(cfg.Netatmo.ClientSecret) == 0 {
		return cfg, errors.New("need a NetAtmo client secret")
	}

	if len(cfg.Netatmo.Username) == 0 {
		return cfg, errors.New("username can not be blank")
	}

	if len(cfg.Netatmo.Password) == 0 {
		return cfg, errors.New("password can not be blank")
	}

	return cfg, nil
}

func main() {
	cfg, err := parseConfig()
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

	http.Handle("/metrics", prometheus.UninstrumentedHandler())
	http.Handle("/", http.RedirectHandler("/metrics", http.StatusFound))

	log.Printf("Listen on %s...", cfg.Addr)
	log.Fatal(http.ListenAndServe(cfg.Addr, nil))
}
