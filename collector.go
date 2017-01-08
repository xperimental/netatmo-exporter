package main

import (
	"log"
	"time"

	netatmo "github.com/exzz/netatmo-api-go"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	staleDataThreshold = 30 * time.Minute
)

var (
	netatmoUp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "netatmo_up",
		Help: "Zero if there was an error scraping the Netatmo API.",
	})

	varLabels = []string{
		"module",
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
		"Atmospheric pressure measurement in mb",
		varLabels,
		nil)
)

type netatmoCollector struct {
	client *netatmo.Client
}

func (m *netatmoCollector) Describe(dChan chan<- *prometheus.Desc) {
	dChan <- updatedDesc
	dChan <- tempDesc
	dChan <- humidityDesc
	dChan <- cotwoDesc
}

func (m *netatmoCollector) Collect(mChan chan<- prometheus.Metric) {
	devices, err := m.client.Read()
	if err != nil {
		netatmoUp.Set(0)
		mChan <- netatmoUp
		return
	}
	netatmoUp.Set(1)
	mChan <- netatmoUp

	for _, dev := range devices.Devices() {
		collectData(mChan, dev)

		for _, module := range dev.LinkedModules {
			collectData(mChan, module)
		}
	}
}

func collectData(ch chan<- prometheus.Metric, device *netatmo.Device) {
	moduleName := device.ModuleName
	data := device.DashboardData

	if data.LastMesure == nil {
		return
	}

	date := time.Unix(*data.LastMesure, 0)
	if time.Since(date) > staleDataThreshold {
		return
	}

	sendMetric(ch, updatedDesc, prometheus.CounterValue, float64(date.UTC().Unix()), moduleName)

	if data.Temperature != nil {
		sendMetric(ch, tempDesc, prometheus.GaugeValue, float64(*data.Temperature), moduleName)
	}

	if data.Humidity != nil {
		sendMetric(ch, humidityDesc, prometheus.GaugeValue, float64(*data.Humidity), moduleName)
	}

	if data.CO2 != nil {
		sendMetric(ch, cotwoDesc, prometheus.GaugeValue, float64(*data.CO2), moduleName)
	}

	if data.Noise != nil {
		sendMetric(ch, noiseDesc, prometheus.GaugeValue, float64(*data.Noise), moduleName)
	}

	if data.Pressure != nil {
		sendMetric(ch, pressureDesc, prometheus.GaugeValue, float64(*data.Pressure), moduleName)
	}
}

func sendMetric(ch chan<- prometheus.Metric, desc *prometheus.Desc, valueType prometheus.ValueType, value float64, moduleName string) {
	m, err := prometheus.NewConstMetric(desc, valueType, value, moduleName)
	if err != nil {
		log.Printf("Error creating %s metric: %s", updatedDesc.String(), err)
	}
	ch <- m
}
