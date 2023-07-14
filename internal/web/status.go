package web

import (
	"fmt"
	"github.com/exzz/netatmo-api-go"
	"golang.org/x/oauth2"
	"html/template"
	"net/http"
	"time"

	_ "embed"
)

var (
	//go:embed home-unauthorized.html
	homeNotAuthorizedHtml string

	//go:embed home-authorized.html
	homeAuthorizedHtml string
)

func HomeHandler(tokenFunc func() (*oauth2.Token, error)) http.Handler {
	homeAuthorizedTemplate, err := template.New("home-authorized.html").Parse(homeAuthorizedHtml)
	if err != nil {
		panic(err)
	}

	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		token, err := tokenFunc()
		switch {
		case err == netatmo.ErrNotAuthenticated:
		case err != nil:
			http.Error(wr, fmt.Sprintf("Error getting token: %s", err), http.StatusInternalServerError)
			return
		default:
		}

		wr.Header().Set("Content-Type", "text/html")

		if !token.Valid() {
			fmt.Fprint(wr, homeNotAuthorizedHtml)
			return
		}

		context := struct {
			Expiry            time.Time
			RemainingDuration time.Duration
		}{
			Expiry:            token.Expiry,
			RemainingDuration: time.Until(token.Expiry),
		}

		if err := homeAuthorizedTemplate.Execute(wr, context); err != nil {
			http.Error(wr, fmt.Sprintf("Error executing template: %s", err), http.StatusInternalServerError)
		}
	})
}
