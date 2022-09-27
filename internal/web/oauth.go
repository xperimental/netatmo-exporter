package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/exzz/netatmo-api-go"
)

func AuthorizeHandler(externalURL string, client *netatmo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		redirectURL := externalURL + "/callback"
		authURL := client.AuthCodeURL(redirectURL, "definitelyrandom")

		http.Redirect(w, r, authURL, http.StatusFound)
	}
}

func CallbackHandler(client *netatmo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		if err := doCallback(r.Context(), client, values); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error processing code: %s", err)
			return
		}

		fmt.Fprintln(w, "Authenticated.")
	}
}

func doCallback(ctx context.Context, client *netatmo.Client, query url.Values) error {
	if err := query.Get("error"); err != "" {
		return errors.New("user did not accept")
	}

	state := query.Get("state")
	code := query.Get("code")

	return client.Exchange(ctx, code, state)
}
