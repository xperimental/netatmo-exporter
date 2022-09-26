package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/exzz/netatmo-api-go"
	"github.com/xperimental/netatmo-exporter/internal/logger"
)

var log = logger.NewLogger()

func main() {
	cfg, err := parseArgs(os.Args[0], os.Args[1:])
	if err != nil {
		log.Fatalf("Error in configuration: %s", err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatalf("Error creating local listener: %s", err)
	}
	defer listener.Close()

	netatmoConfig := netatmo.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
	}
	client := netatmo.NewClient(netatmoConfig)

	redirectURL := fmt.Sprintf("http://%s/callback", listener.Addr())

	mux := http.NewServeMux()
	mux.Handle("/", indexHandler(client))
	mux.Handle("/authorize", authorizeHandler(client, redirectURL))
	mux.Handle("/callback", callbackHandler(client))
	mux.Handle("/token", tokenHandler(client))

	log.Infof("Open this link in your browser: http://%s", listener.Addr())

	server := &http.Server{
		Handler: mux,
	}
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
