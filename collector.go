package main

import (
	"time"

	netatmo "github.com/exzz/netatmo-api-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

var (
	netatmoUp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "netatmo_up",
		Help: "Zero if there was an error scraping the Netatmo API.",
	})

	varLabels = []string{
		"module",
		"station",
	}

	prefix = "netatmo_sensor_"

	updatedDesc = prometheus.NewDesc(
		prefix+"updated",
		"Timestamp of last update",
		varLabels,
		nil)

	tempDesc = prometheus.NewDesc(
		prefix+"temperature_celsius",
		"Temperature measurement in celsius",
		varLabels,
		nil)

	humidityDesc = prometheus.NewDesc(
		prefix+"humidity_percent",
		"Relative humidity measurement in percent",
		varLabels,
		nil)

	cotwoDesc = prometheus.NewDesc(
		prefix+"co2_ppm",
		"Carbondioxide measurement in parts per million",
		varLabels,
		nil)

	noiseDesc = prometheus.NewDesc(
		prefix+"noise_db",
		"Noise measurement in decibels",
		varLabels,
		nil)

	pressureDesc = prometheus.NewDesc(
		prefix+"pressure_mb",
		"Atmospheric pressure measurement in millibar",
		varLabels,
		nil)

	windStrengthDesc = prometheus.NewDesc(
		prefix+"wind_strength_kph",
		"Wind strength in kilometers per hour",
		varLabels,
		nil)

	windDirectionDesc = prometheus.NewDesc(
		prefix+"wind_direction_degrees",
		"Wind direction in degrees",
		varLabels,
		nil)

	rainDesc = prometheus.NewDesc(
		prefix+"rain_amount_mm",
		"Rain amount in millimeters",
		varLabels,
		nil)

	batteryDesc = prometheus.NewDesc(
		prefix+"battery_percent",
		"Battery remaining life (10: low)",
		varLabels,
		nil)
	wifiDesc = prometheus.NewDesc(
		prefix+"wifi_signal_strength",
		"Wifi signal strength (86: bad, 71: avg, 56: good)",
		varLabels,
		nil)
	rfDesc = prometheus.NewDesc(
		prefix+"rf_signal_strength",
		"RF signal strength (90: lowest, 60: highest)",
		varLabels,
		nil)
)

type netatmoCollector struct {
	log            logrus.FieldLogger
	staleThreshold time.Duration
	client         *netatmo.Client
}

func (c *netatmoCollector) Describe(dChan chan<- *prometheus.Desc) {
	dChan <- updatedDesc
	dChan <- tempDesc
	dChan <- humidityDesc
	dChan <- cotwoDesc
}

func (c *netatmoCollector) Collect(mChan chan<- prometheus.Metric) {
	devices, err := c.client.Read()
	if err != nil {
		c.log.Errorf("Error getting data: %s", err)

		netatmoUp.Set(0)
		mChan <- netatmoUp
		return
	}
	netatmoUp.Set(1)
	mChan <- netatmoUp

	for _, dev := range devices.Devices() {
		stationName := dev.StationName
		c.collectData(mChan, dev, stationName)

		for _, module := range dev.LinkedModules {
			c.collectData(mChan, module, stationName)
		}
	}
}

func (c *netatmoCollector) collectData(ch chan<- prometheus.Metric, device *netatmo.Device, stationName string) {
	moduleName := device.ModuleName
	data := device.DashboardData

	if data.LastMeasure == nil {
		c.log.Debugf("No data available.")
		return
	}

	date := time.Unix(*data.LastMeasure, 0)
	if time.Since(date) > c.staleThreshold {
		c.log.Warnf("Data is stale for %s: %s > %s", moduleName, time.Since(date), c.staleThreshold)
		return
	}

	c.sendMetric(ch, updatedDesc, prometheus.CounterValue, float64(date.UTC().Unix()), moduleName, stationName)

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

func (c *netatmoCollector) sendMetric(ch chan<- prometheus.Metric, desc *prometheus.Desc, valueType prometheus.ValueType, value float64, moduleName string, stationName string) {
	m, err := prometheus.NewConstMetric(desc, valueType, value, moduleName, stationName)
	if err != nil {
		c.log.Errorf("Error creating %s metric: %s", updatedDesc.String(), err)
	}
	ch <- m
}
