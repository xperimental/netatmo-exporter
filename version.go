package main

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

var (
	// Version contains the version as set during the build.
	Version = "unknown"

	// GitCommit contains the git commit hash set during the build.
	GitCommit = "unknown"
)

func versionHandler(log logrus.FieldLogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		info := struct {
			Version string `json:"version"`
			Commit  string `json:"commit"`
		}{
			Version: Version,
			Commit:  GitCommit,
		}

		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(info); err != nil {
			log.Errorf("Error encoding version info: %s", err)
		}
	})
}
