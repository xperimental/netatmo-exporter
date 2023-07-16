package token

import (
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/oauth2"
)

const (
	prefix = "netatmo_exporter_token_"
)

var (
	validDesc = prometheus.NewDesc(
		prefix+"valid",
		"Set to 1 if there is a valid token, 0 otherwise.",
		nil, nil)

	expiryDesc = prometheus.NewDesc(
		prefix+"expiry_time",
		"Set to the unix timestamp when the token will expire. 0 if no expiry is set.",
		nil, nil)
)

func Metric(tokenFunc func() (*oauth2.Token, error)) prometheus.Collector {
	return &tokenMetric{
		tokenFunc: tokenFunc,
	}
}

type tokenMetric struct {
	tokenFunc func() (*oauth2.Token, error)
}

func (t tokenMetric) Describe(dChan chan<- *prometheus.Desc) {
	dChan <- validDesc
	dChan <- expiryDesc
}

func (t tokenMetric) Collect(mChan chan<- prometheus.Metric) {
	token, _ := t.tokenFunc()

	valid := token.Valid()
	validValue := 0.0
	expiryValue := 0.0
	if valid {
		validValue = 1.0
		expiryValue = float64(token.Expiry.Unix())
	}

	mChan <- prometheus.MustNewConstMetric(validDesc, prometheus.GaugeValue, validValue)
	mChan <- prometheus.MustNewConstMetric(expiryDesc, prometheus.GaugeValue, expiryValue)
}
