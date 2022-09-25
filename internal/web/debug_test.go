package web

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/exzz/netatmo-api-go"
	"github.com/google/go-cmp/cmp"
	"github.com/sirupsen/logrus"
)

func TestDebugHandler(t *testing.T) {
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
			desc: "success",
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
			h := DebugHandler(log, tc.readFunc)

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
