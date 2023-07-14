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
	//go:embed home.html
	homeHtml string
)

func HomeHandler(tokenFunc func() (*oauth2.Token, error)) http.Handler {
	homeTemplate, err := template.New("home.html").Funcs(map[string]any{
		"remaining": remaining,
	}).Parse(homeHtml)
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

		context := struct {
			Valid bool
			Token *oauth2.Token
		}{
			Valid: token.Valid(),
			Token: token,
		}

		wr.Header().Set("Content-Type", "text/html")
		if err := homeTemplate.Execute(wr, context); err != nil {
			http.Error(wr, fmt.Sprintf("Error executing template: %s", err), http.StatusInternalServerError)
		}
	})
}

func remaining(t time.Time) time.Duration {
	return time.Until(t)
}
