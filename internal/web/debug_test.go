package web

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/exzz/netatmo-api-go"
	"github.com/google/go-cmp/cmp"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func TestDebugDataHandler(t *testing.T) {
	createCollection := func(devices []*netatmo.Device) *netatmo.DeviceCollection {
		dc := &netatmo.DeviceCollection{}
		dc.Body.Devices = devices
		return dc
	}
	tt := []struct {
		desc       string
		readFunc   func() (*netatmo.DeviceCollection, error)
		wantStatus int
		wantBody   string
	}{
		{
			desc: "success",
			readFunc: func() (*netatmo.DeviceCollection, error) {
				return createCollection([]*netatmo.Device{}), nil
			},
			wantStatus: http.StatusOK,
			wantBody: `{"Body":{"devices":[]}}
`,
		},
		{
			desc: "error retrieving data",
			readFunc: func() (*netatmo.DeviceCollection, error) {
				return nil, errors.New("test error")
			},
			wantStatus: http.StatusBadGateway,
			wantBody: `Error retrieving data: test error
`,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)

			log := logrus.New()
			h := DebugDataHandler(log, tc.readFunc)

			h.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("got code %d, want %d", rec.Code, tc.wantStatus)
			}

			body := rec.Body.String()
			if diff := cmp.Diff(body, tc.wantBody); diff != "" {
				t.Errorf("body differs: -got+want\n%s", diff)
			}
		})
	}
}

func TestDebugTokenHandler(t *testing.T) {
	tt := []struct {
		desc       string
		tokenFunc  func() (*oauth2.Token, error)
		wantStatus int
		wantBody   string
	}{
		{
			desc: "success",
			tokenFunc: func() (*oauth2.Token, error) {
				return &oauth2.Token{
					AccessToken:  "access-token",
					RefreshToken: "refresh-token",
					Expiry:       time.Unix(0, 0).UTC(),
				}, nil
			},
			wantStatus: http.StatusOK,
			wantBody: `{
  "isValid": false,
  "hasAccessToken": true,
  "hasRefreshToken": true,
  "expiry": "1970-01-01T00:00:00Z"
}
`,
		},
		{
			desc: "no token",
			tokenFunc: func() (*oauth2.Token, error) {
				return nil, netatmo.ErrNotAuthenticated
			},
			wantStatus: http.StatusNotFound,
			wantBody:   "No token available.\n",
		},
		{
			desc: "error retrieving token",
			tokenFunc: func() (*oauth2.Token, error) {
				return nil, errors.New("error")
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "Error retrieving token: error\n",
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)

			log := logrus.New()
			h := DebugTokenHandler(log, tc.tokenFunc)

			h.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("got code %d, want %d", rec.Code, tc.wantStatus)
			}

			body := rec.Body.String()
			if diff := cmp.Diff(body, tc.wantBody); diff != "" {
				t.Errorf("body differs: -got+want\n%s", diff)
			}
		})
	}
}
