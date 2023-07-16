package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/exzz/netatmo-api-go"
	"golang.org/x/oauth2"
)

func AuthorizeHandler(externalURL string, client *netatmo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		redirectURL := externalURL + "/auth/callback"
		authURL := client.AuthCodeURL(redirectURL, "definitelyrandom")

		http.Redirect(w, r, authURL, http.StatusFound)
	}
}

func CallbackHandler(ctx context.Context, client *netatmo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		if err := doCallback(ctx, client, values); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error processing code: %s", err)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
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

func SetTokenHandler(ctx context.Context, client *netatmo.Client) http.HandlerFunc {
	return func(wr http.ResponseWriter, r *http.Request) {
		refreshToken := r.FormValue("refresh_token")
		if refreshToken == "" {
			http.Error(wr, "The refresh token can not be empty. Please go back.", http.StatusBadRequest)
			return
		}

		token := &oauth2.Token{
			RefreshToken: refreshToken,
		}
		client.InitWithToken(ctx, token)

		http.Redirect(wr, r, "/", http.StatusFound)
	}
}
