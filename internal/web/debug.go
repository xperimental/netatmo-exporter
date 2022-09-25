package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/exzz/netatmo-api-go"
	"github.com/sirupsen/logrus"
)

// DebugHandler creates a handler which outputs the raw JSON data.
func DebugHandler(log logrus.FieldLogger, readFunc func() (*netatmo.DeviceCollection, error)) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		devices, err := readFunc()
		if err != nil {
			http.Error(wr, fmt.Sprintf("Error retrieving data: %s", err), http.StatusBadGateway)
			return
		}

		wr.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(wr).Encode(devices); err != nil {
			log.Errorf("Can not encode response: %s", err)
			return
		}
	})
}
