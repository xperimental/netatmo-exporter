package collector

import (
	"errors"
	"testing"
	"time"

	netatmo "github.com/exzz/netatmo-api-go"
	"github.com/google/go-cmp/cmp"
	"github.com/sirupsen/logrus"
)

func TestRefreshData(t *testing.T) {
	testData := &netatmo.DeviceCollection{}
	testError := errors.New("test error")
	tt := []struct {
		desc         string
		time         time.Time
		readFunction ReadFunction
		wantTime     time.Time
		wantData     *netatmo.DeviceCollection
		wantError    error
	}{
		{
			desc: "success",
			time: time.Unix(0, 0),
			readFunction: func() (*netatmo.DeviceCollection, error) {
				return testData, nil
			},
			wantTime:  time.Unix(0, 0),
			wantData:  testData,
			wantError: nil,
		},
		{
			desc: "error",
			time: time.Unix(0, 0),
			readFunction: func() (*netatmo.DeviceCollection, error) {
				return nil, testError
			},
			wantTime:  time.Time{},
			wantData:  nil,
			wantError: testError,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			c := &NetatmoCollector{
				Log:          logrus.New(),
				ReadFunction: tc.readFunction,
			}
			c.RefreshData(tc.time)

			if c.cacheTimestamp != tc.wantTime {
				t.Errorf("got time %s, want %s", c.cacheTimestamp, tc.wantTime)
			}

			if diff := cmp.Diff(c.cachedData, tc.wantData); diff != "" {
				t.Errorf("data differs: -got+want\n%s", diff)
			}

			if c.lastRefreshError != tc.wantError {
				t.Errorf("got error %q, want %q", c.lastRefreshError, tc.wantError)
			}
		})
	}
}

func TestRefreshDataResetError(t *testing.T) {
	testData := &netatmo.DeviceCollection{}
	testError := errors.New("test error")
	successFunc := func() (*netatmo.DeviceCollection, error) {
		return testData, nil
	}
	errorFunc := func() (*netatmo.DeviceCollection, error) {
		return nil, testError
	}

	c := &NetatmoCollector{
		Log:          logrus.New(),
		ReadFunction: successFunc,
	}
	c.RefreshData(time.Unix(0, 0))

	if c.lastRefreshError != nil {
		t.Errorf("got error %q, want none", c.lastRefreshError)
	}

	c.ReadFunction = errorFunc
	c.RefreshData(time.Unix(1, 0))

	if c.lastRefreshError != testError {
		t.Errorf("got error %q, want %q", c.lastRefreshError, testError)
	}

	c.ReadFunction = successFunc
	c.RefreshData(time.Unix(0, 0))

	if c.lastRefreshError != nil {
		t.Errorf("got error %q, want none", c.lastRefreshError)
	}
}
