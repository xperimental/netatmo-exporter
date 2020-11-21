package collector

import (
	"sync"
	"time"

	netatmo "github.com/exzz/netatmo-api-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

var (
	prefix        = "netatmo_"
	netatmoUpDesc = prometheus.NewDesc(prefix+"up",
		"Zero if there was an error during the last refresh try.",
		nil, nil)

	refreshIntervalDesc = prometheus.NewDesc(
		prefix+"refresh_interval_seconds",
		"Contains the configured refresh interval in seconds. This is provided as a convenience for calculations with the cache update time.",
		nil, nil)
	refreshPrefix        = prefix + "last_refresh_"
	refreshTimestampDesc = prometheus.NewDesc(
		refreshPrefix+"time",
		"Contains the time of the last refresh try, successful or not.",
		nil, nil)
	refreshDurationDesc = prometheus.NewDesc(
		refreshPrefix+"duration_seconds",
		"Contains the time it took for the last refresh to complete, even if it was unsuccessful.",
		nil, nil)

	cacheTimestampDesc = prometheus.NewDesc(
		prefix+"cache_updated_time",
		"Contains the time of the cached data.",
		nil, nil)

	varLabels = []string{
		"module",
		"station",
	}

	sensorPrefix = prefix + "sensor_"

	updatedDesc = prometheus.NewDesc(
		sensorPrefix+"updated",
		"Timestamp of last update",
		varLabels,
		nil)

	tempDesc = prometheus.NewDesc(
		sensorPrefix+"temperature_celsius",
		"Temperature measurement in celsius",
		varLabels,
		nil)

	humidityDesc = prometheus.NewDesc(
		sensorPrefix+"humidity_percent",
		"Relative humidity measurement in percent",
		varLabels,
		nil)

	cotwoDesc = prometheus.NewDesc(
		sensorPrefix+"co2_ppm",
		"Carbondioxide measurement in parts per million",
		varLabels,
		nil)

	noiseDesc = prometheus.NewDesc(
		sensorPrefix+"noise_db",
		"Noise measurement in decibels",
		varLabels,
		nil)

	pressureDesc = prometheus.NewDesc(
		sensorPrefix+"pressure_mb",
		"Atmospheric pressure measurement in millibar",
		varLabels,
		nil)

	windStrengthDesc = prometheus.NewDesc(
		sensorPrefix+"wind_strength_kph",
		"Wind strength in kilometers per hour",
		varLabels,
		nil)

	windDirectionDesc = prometheus.NewDesc(
		sensorPrefix+"wind_direction_degrees",
		"Wind direction in degrees",
		varLabels,
		nil)

	rainDesc = prometheus.NewDesc(
		sensorPrefix+"rain_amount_mm",
		"Rain amount in millimeters",
		varLabels,
		nil)

	batteryDesc = prometheus.NewDesc(
		sensorPrefix+"battery_percent",
		"Battery remaining life (10: low)",
		varLabels,
		nil)
	wifiDesc = prometheus.NewDesc(
		sensorPrefix+"wifi_signal_strength",
		"Wifi signal strength (86: bad, 71: avg, 56: good)",
		varLabels,
		nil)
	rfDesc = prometheus.NewDesc(
		sensorPrefix+"rf_signal_strength",
		"RF signal strength (90: lowest, 60: highest)",
		varLabels,
		nil)
)

// ReadFunction defines the interface for reading from the Netatmo API.
type ReadFunction func() (*netatmo.DeviceCollection, error)

// NetatmoCollector is a Prometheus collector for Netatmo sensor values.
type NetatmoCollector struct {
	Log                 logrus.FieldLogger
	RefreshInterval     time.Duration
	StaleThreshold      time.Duration
	ReadFunction        ReadFunction
	lastRefresh         time.Time
	lastRefreshError    error
	lastRefreshDuration time.Duration
	cacheLock           sync.RWMutex
	cacheTimestamp      time.Time
	cachedData          *netatmo.DeviceCollection
}

// Describe implements prometheus.Collector
func (c *NetatmoCollector) Describe(dChan chan<- *prometheus.Desc) {
	dChan <- updatedDesc
	dChan <- tempDesc
	dChan <- humidityDesc
	dChan <- cotwoDesc
}

// Collect implements prometheus.Collector
func (c *NetatmoCollector) Collect(mChan chan<- prometheus.Metric) {
	now := time.Now()
	if now.Sub(c.lastRefresh) >= c.RefreshInterval {
		go c.RefreshData(now)
	}

	upValue := 1.0
	if c.lastRefresh.IsZero() || c.lastRefreshError != nil {
		upValue = 0
	}
	c.sendMetric(mChan, netatmoUpDesc, prometheus.GaugeValue, upValue)
	c.sendMetric(mChan, refreshIntervalDesc, prometheus.GaugeValue, c.RefreshInterval.Seconds())
	c.sendMetric(mChan, refreshTimestampDesc, prometheus.GaugeValue, convertTime(c.lastRefresh))
	c.sendMetric(mChan, refreshDurationDesc, prometheus.GaugeValue, c.lastRefreshDuration.Seconds())

	c.cacheLock.RLock()
	defer c.cacheLock.RUnlock()

	c.sendMetric(mChan, cacheTimestampDesc, prometheus.GaugeValue, convertTime(c.cacheTimestamp))
	if c.cachedData != nil {
		for _, dev := range c.cachedData.Devices() {
			stationName := dev.StationName
			c.collectData(mChan, dev, stationName)

			for _, module := range dev.LinkedModules {
				c.collectData(mChan, module, stationName)
			}
		}
	}
}

// RefreshData causes the collector to try to refresh the cached data.
func (c *NetatmoCollector) RefreshData(now time.Time) {
	c.Log.Debugf("Refreshing data. Time since last refresh: %s", now.Sub(c.lastRefresh))
	c.lastRefresh = now

	defer func(start time.Time) {
		c.lastRefreshDuration = time.Since(start)
	}(time.Now())

	devices, err := c.ReadFunction()
	c.lastRefreshError = err
	if err != nil {
		c.Log.Errorf("Error during refresh: %s", err)
		return
	}

	c.cacheLock.Lock()
	defer c.cacheLock.Unlock()
	c.cacheTimestamp = now
	c.cachedData = devices
}

func (c *NetatmoCollector) collectData(ch chan<- prometheus.Metric, device *netatmo.Device, stationName string) {
	moduleName := device.ModuleName
	data := device.DashboardData

	if data.LastMeasure == nil {
		c.Log.Debugf("No data available.")
		return
	}

	date := time.Unix(*data.LastMeasure, 0)
	if time.Since(date) > c.StaleThreshold {
		c.Log.Debugf("Data is stale for %s: %s > %s", moduleName, time.Since(date), c.StaleThreshold)
		return
	}

	c.sendMetric(ch, updatedDesc, prometheus.GaugeValue, float64(date.UTC().Unix()), moduleName, stationName)

	if data.Temperature != nil {
		c.sendMetric(ch, tempDesc, prometheus.GaugeValue, float64(*data.Temperature), moduleName, stationName)
	}

	if data.Humidity != nil {
		c.sendMetric(ch, humidityDesc, prometheus.GaugeValue, float64(*data.Humidity), moduleName, stationName)
	}

	if data.CO2 != nil {
		c.sendMetric(ch, cotwoDesc, prometheus.GaugeValue, float64(*data.CO2), moduleName, stationName)
	}

	if data.Noise != nil {
		c.sendMetric(ch, noiseDesc, prometheus.GaugeValue, float64(*data.Noise), moduleName, stationName)
	}

	if data.Pressure != nil {
		c.sendMetric(ch, pressureDesc, prometheus.GaugeValue, float64(*data.Pressure), moduleName, stationName)
	}

	if data.WindStrength != nil {
		c.sendMetric(ch, windStrengthDesc, prometheus.GaugeValue, float64(*data.WindStrength), moduleName, stationName)
	}

	if data.WindAngle != nil {
		c.sendMetric(ch, windDirectionDesc, prometheus.GaugeValue, float64(*data.WindAngle), moduleName, stationName)
	}

	if data.Rain != nil {
		c.sendMetric(ch, rainDesc, prometheus.GaugeValue, float64(*data.Rain), moduleName, stationName)
	}

	if device.BatteryPercent != nil {
		c.sendMetric(ch, batteryDesc, prometheus.GaugeValue, float64(*device.BatteryPercent), moduleName, stationName)
	}
	if device.WifiStatus != nil {
		c.sendMetric(ch, wifiDesc, prometheus.GaugeValue, float64(*device.WifiStatus), moduleName, stationName)
	}
	if device.RFStatus != nil {
		c.sendMetric(ch, rfDesc, prometheus.GaugeValue, float64(*device.RFStatus), moduleName, stationName)
	}
}

func (c *NetatmoCollector) sendMetric(ch chan<- prometheus.Metric, desc *prometheus.Desc, valueType prometheus.ValueType, value float64, labelValues ...string) {
	m, err := prometheus.NewConstMetric(desc, valueType, value, labelValues...)
	if err != nil {
		c.Log.Errorf("Error creating %s metric: %s", updatedDesc.String(), err)
		return
	}
	ch <- m
}

func convertTime(t time.Time) float64 {
	if t.IsZero() {
		return 0.0
	}

	return float64(t.Unix())
}
