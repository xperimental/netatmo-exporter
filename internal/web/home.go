package web

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/exzz/netatmo-api-go"
	"golang.org/x/oauth2"

	_ "embed"
)

const netatmoDevSite = "https://dev.netatmo.com/apps/"

//go:embed home.html
var homeHtml string

type homeContext struct {
	Valid          bool
	Token          *oauth2.Token
	NetAtmoDevSite string
}

// HomeHandler produces a simple website showing the exporter's status in a human-readable form.
// It provides links to other information and help for authentication as well.
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

		context := homeContext{
			Valid:          token.Valid(),
			Token:          token,
			NetAtmoDevSite: netatmoDevSite,
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
