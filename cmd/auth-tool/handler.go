package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/exzz/netatmo-api-go"
	"golang.org/x/oauth2"
)

const (
	veryRandomState = "definitelyrandom"
)

var (
	//go:embed _assets
	assetFS embed.FS

	templates *template.Template
)

func init() {
	assets, err := fs.Sub(assetFS, "_assets")
	if err != nil {
		panic(err)
	}

	templates = template.Must(template.New("").ParseFS(assets, "*.html"))
}

type indexContext struct {
	Token *oauth2.Token
}

func indexHandler(client *netatmo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := client.CurrentToken()
		if err != nil {
			token = nil
		}

		data := indexContext{
			Token: token,
		}

		buf := &bytes.Buffer{}
		if err := templates.ExecuteTemplate(buf, "index.html", data); err != nil {
			http.Error(w, fmt.Sprintf("Error executing template: %s", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		if _, err := w.Write(buf.Bytes()); err != nil {
			log.Errorf("Can not write response: %s", err)
		}
	}
}

func authorizeHandler(client *netatmo.Client, redirectURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authURL := client.AuthCodeURL(redirectURL, veryRandomState)
		http.Redirect(w, r, authURL, http.StatusFound)
	}
}

func callbackHandler(client *netatmo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		if err := client.Exchange(r.Context(), query.Get("code"), query.Get("state")); err != nil {
			http.Error(w, fmt.Sprintf("Error during authorization: %s", err), http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func tokenHandler(client *netatmo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := client.CurrentToken()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading token: %s", err), http.StatusInternalServerError)
			return
		}

		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetIndent("", "  ")
		if err := enc.Encode(token); err != nil {
			http.Error(w, fmt.Sprintf("Error marshalling token: %s", err), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Disposition", `attachment; filename="token.json"`)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(buf.Bytes()); err != nil {
			log.Errorf("Can not write response: %s", err)
		}
	}
}
