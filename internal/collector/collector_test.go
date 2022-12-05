package collector

import (
	"errors"
	"strings"
	"testing"
	"time"

	netatmo "github.com/exzz/netatmo-api-go"
	"github.com/google/go-cmp/cmp"
	"github.com/prometheus/client_golang/prometheus/testutil"
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

			c := New(logrus.New(), tc.readFunction, 0, 0)
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

	c := New(logrus.New(), successFunc, 0, 0)
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

func TestNetatmoCollector_Collect(t *testing.T) {
	tt := []struct {
		desc        string
		data        *netatmo.DeviceCollection
		wantMetrics string
	}{
		{
			desc: "success, no data",
			data: &netatmo.DeviceCollection{},
			wantMetrics: `# HELP netatmo_cache_updated_time Contains the time of the cached data.
# TYPE netatmo_cache_updated_time gauge
netatmo_cache_updated_time 1
# HELP netatmo_last_refresh_duration_seconds Contains the time it took for the last refresh to complete, even if it was unsuccessful.
# TYPE netatmo_last_refresh_duration_seconds gauge
netatmo_last_refresh_duration_seconds 0
# HELP netatmo_last_refresh_time Contains the time of the last refresh try, successful or not.
# TYPE netatmo_last_refresh_time gauge
netatmo_last_refresh_time 1
# HELP netatmo_refresh_interval_seconds Contains the configured refresh interval in seconds. This is provided as a convenience for calculations with the cache update time.
# TYPE netatmo_refresh_interval_seconds gauge
netatmo_refresh_interval_seconds 3600
# HELP netatmo_up Zero if there was an error during the last refresh try.
# TYPE netatmo_up gauge
netatmo_up 1
`,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			mockClock := func() time.Time {
				return time.Unix(1, 0)
			}

			read := func() (*netatmo.DeviceCollection, error) {
				return tc.data, nil
			}
			expected := strings.NewReader(tc.wantMetrics)

			c := New(logrus.New(), read, time.Hour, time.Hour)
			c.clock = mockClock
			c.RefreshData(mockClock())

			if err := testutil.CollectAndCompare(c, expected); err != nil {
				t.Error(err)
			}
		})
	}
}
