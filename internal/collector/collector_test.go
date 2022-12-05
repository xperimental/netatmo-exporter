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
	testDevices := &netatmo.DeviceCollection{}
	testDevices.Body.Devices = []*netatmo.Device{
		{
			ID:          "aa:bb:cc:dd:ee:f0",
			ModuleName:  "Living Room",
			HomeID:      "0123456789abcdef01234567",
			HomeName:    "Home",
			StationName: "Home (Living Room)",
			WifiStatus:  int32Ptr(45),
			Type:        "NAMain",
			DashboardData: netatmo.DashboardData{
				Temperature:      float32Ptr(23),
				Humidity:         int32Ptr(45),
				CO2:              int32Ptr(650),
				Noise:            int32Ptr(40),
				Pressure:         float32Ptr(1234),
				AbsolutePressure: float32Ptr(987),
				LastMeasure:      int64Ptr(3500),
			},
			LinkedModules: []*netatmo.Device{
				{
					ID:             "aa:bb:cc:dd:ee:f1",
					ModuleName:     "Outside",
					BatteryPercent: int32Ptr(70),
					RFStatus:       int32Ptr(57),
					Type:           "NAModule1",
					DashboardData: netatmo.DashboardData{
						Temperature: float32Ptr(5),
						Humidity:    int32Ptr(83),
						LastMeasure: int64Ptr(3501),
					},
				},
				{
					ID:             "aa:bb:cc:dd:ee:f2",
					ModuleName:     "Bedroom",
					BatteryPercent: int32Ptr(55),
					RFStatus:       int32Ptr(80),
					Type:           "NAModule4",
					DashboardData: netatmo.DashboardData{
						Temperature: float32Ptr(17),
						Humidity:    int32Ptr(52),
						CO2:         int32Ptr(510),
						LastMeasure: int64Ptr(3502),
					},
				},
			},
		},
	}

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
		netatmo_cache_updated_time 3600
		# HELP netatmo_last_refresh_duration_seconds Contains the time it took for the last refresh to complete, even if it was unsuccessful.
		# TYPE netatmo_last_refresh_duration_seconds gauge
		netatmo_last_refresh_duration_seconds 0
		# HELP netatmo_last_refresh_time Contains the time of the last refresh try, successful or not.
		# TYPE netatmo_last_refresh_time gauge
		netatmo_last_refresh_time 3600
		# HELP netatmo_refresh_interval_seconds Contains the configured refresh interval in seconds. This is provided as a convenience for calculations with the cache update time.
		# TYPE netatmo_refresh_interval_seconds gauge
		netatmo_refresh_interval_seconds 3600
		# HELP netatmo_up Zero if there was an error during the last refresh try.
		# TYPE netatmo_up gauge
		netatmo_up 1
		`,
		},
		{
			desc: "success",
			data: testDevices,
			wantMetrics: `# HELP netatmo_cache_updated_time Contains the time of the cached data.
# TYPE netatmo_cache_updated_time gauge
netatmo_cache_updated_time 3600
# HELP netatmo_last_refresh_duration_seconds Contains the time it took for the last refresh to complete, even if it was unsuccessful.
# TYPE netatmo_last_refresh_duration_seconds gauge
netatmo_last_refresh_duration_seconds 0
# HELP netatmo_last_refresh_time Contains the time of the last refresh try, successful or not.
# TYPE netatmo_last_refresh_time gauge
netatmo_last_refresh_time 3600
# HELP netatmo_refresh_interval_seconds Contains the configured refresh interval in seconds. This is provided as a convenience for calculations with the cache update time.
# TYPE netatmo_refresh_interval_seconds gauge
netatmo_refresh_interval_seconds 3600
# HELP netatmo_sensor_battery_percent Battery remaining life (10: low)
# TYPE netatmo_sensor_battery_percent gauge
netatmo_sensor_battery_percent{home="Home",module="Bedroom",station="Home (Living Room)"} 55
netatmo_sensor_battery_percent{home="Home",module="Outside",station="Home (Living Room)"} 70
# HELP netatmo_sensor_co2_ppm Carbondioxide measurement in parts per million
# TYPE netatmo_sensor_co2_ppm gauge
netatmo_sensor_co2_ppm{home="Home",module="Bedroom",station="Home (Living Room)"} 510
netatmo_sensor_co2_ppm{home="Home",module="Living Room",station="Home (Living Room)"} 650
# HELP netatmo_sensor_humidity_percent Relative humidity measurement in percent
# TYPE netatmo_sensor_humidity_percent gauge
netatmo_sensor_humidity_percent{home="Home",module="Bedroom",station="Home (Living Room)"} 52
netatmo_sensor_humidity_percent{home="Home",module="Living Room",station="Home (Living Room)"} 45
netatmo_sensor_humidity_percent{home="Home",module="Outside",station="Home (Living Room)"} 83
# HELP netatmo_sensor_noise_db Noise measurement in decibels
# TYPE netatmo_sensor_noise_db gauge
netatmo_sensor_noise_db{home="Home",module="Living Room",station="Home (Living Room)"} 40
# HELP netatmo_sensor_pressure_mb Atmospheric pressure measurement in millibar
# TYPE netatmo_sensor_pressure_mb gauge
netatmo_sensor_pressure_mb{home="Home",module="Living Room",station="Home (Living Room)"} 1234
# HELP netatmo_sensor_rf_signal_strength RF signal strength (90: lowest, 60: highest)
# TYPE netatmo_sensor_rf_signal_strength gauge
netatmo_sensor_rf_signal_strength{home="Home",module="Bedroom",station="Home (Living Room)"} 80
netatmo_sensor_rf_signal_strength{home="Home",module="Outside",station="Home (Living Room)"} 57
# HELP netatmo_sensor_temperature_celsius Temperature measurement in celsius
# TYPE netatmo_sensor_temperature_celsius gauge
netatmo_sensor_temperature_celsius{home="Home",module="Bedroom",station="Home (Living Room)"} 17
netatmo_sensor_temperature_celsius{home="Home",module="Living Room",station="Home (Living Room)"} 23
netatmo_sensor_temperature_celsius{home="Home",module="Outside",station="Home (Living Room)"} 5
# HELP netatmo_sensor_updated Timestamp of last update
# TYPE netatmo_sensor_updated gauge
netatmo_sensor_updated{home="Home",module="Bedroom",station="Home (Living Room)"} 3502
netatmo_sensor_updated{home="Home",module="Living Room",station="Home (Living Room)"} 3500
netatmo_sensor_updated{home="Home",module="Outside",station="Home (Living Room)"} 3501
# HELP netatmo_sensor_wifi_signal_strength Wifi signal strength (86: bad, 71: avg, 56: good)
# TYPE netatmo_sensor_wifi_signal_strength gauge
netatmo_sensor_wifi_signal_strength{home="Home",module="Living Room",station="Home (Living Room)"} 45
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
				return time.Unix(3600, 0)
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

func int32Ptr(i int32) *int32 {
	return &i
}

func int64Ptr(i int64) *int64 {
	return &i
}

func float32Ptr(f float32) *float32 {
	return &f
}
