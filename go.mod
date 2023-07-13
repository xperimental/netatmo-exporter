module github.com/xperimental/netatmo-exporter

go 1.20

require (
	github.com/exzz/netatmo-api-go v0.0.0-20201009073308-a8620474d1ea
	github.com/google/go-cmp v0.5.9
	github.com/prometheus/client_golang v1.16.0
	github.com/prometheus/procfs v0.11.0 // indirect
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/pflag v1.0.5
	golang.org/x/net v0.12.0 // indirect
	golang.org/x/oauth2 v0.10.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/prometheus/client_model v0.4.0 // indirect
	github.com/prometheus/common v0.44.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
)

replace github.com/exzz/netatmo-api-go => github.com/xperimental/netatmo-api-go v0.0.0-20220918192520-fb6577104ce6
