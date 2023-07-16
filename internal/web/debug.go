package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/exzz/netatmo-api-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// DebugDataHandler creates a handler which outputs the raw JSON data.
func DebugDataHandler(log logrus.FieldLogger, readFunc func() (*netatmo.DeviceCollection, error)) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		devices, err := readFunc()
		if err != nil {
			http.Error(wr, fmt.Sprintf("Error retrieving data: %s", err), http.StatusBadGateway)
			return
		}

		wr.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(wr).Encode(devices); err != nil {
			log.Errorf("Can not encode data debug response: %s", err)
			return
		}
	})
}

// DebugTokenHandler creates a handler which returns information about the currently-used token.
// For security reasons, the actual token data is not returned.
func DebugTokenHandler(log logrus.FieldLogger, tokenFunc func() (*oauth2.Token, error)) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		token, err := tokenFunc()
		switch {
		case err == netatmo.ErrNotAuthenticated:
		case err != nil:
			http.Error(wr, fmt.Sprintf("Error retrieving token: %s", err), http.StatusInternalServerError)
			return
		default:
		}

		if token == nil {
			http.Error(wr, "No token available.", http.StatusNotFound)
			return
		}

		data := struct {
			IsValid         bool      `json:"isValid"`
			HasAccessToken  bool      `json:"hasAccessToken"`
			HasRefreshToken bool      `json:"hasRefreshToken"`
			Expiry          time.Time `json:"expiry"`
		}{
			IsValid:         token.Valid(),
			HasAccessToken:  token.AccessToken != "",
			HasRefreshToken: token.RefreshToken != "",
			Expiry:          token.Expiry,
		}

		wr.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(wr)
		enc.SetIndent("", "  ")
		if err := enc.Encode(data); err != nil {
			log.Errorf("Can not encode token debug response: %s", err)
			return
		}
	})
}
