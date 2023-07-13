module github.com/xperimental/netatmo-exporter

go 1.20

require (
	github.com/exzz/netatmo-api-go v0.0.0-20201009073308-a8620474d1ea
	github.com/google/go-cmp v0.5.9
	github.com/prometheus/client_golang v1.14.0
	github.com/prometheus/procfs v0.9.0 // indirect
	github.com/sirupsen/logrus v1.9.0
	github.com/spf13/pflag v1.0.5
	golang.org/x/net v0.5.0 // indirect
	golang.org/x/oauth2 v0.4.0
	golang.org/x/sys v0.4.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.39.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)

replace github.com/exzz/netatmo-api-go => github.com/xperimental/netatmo-api-go v0.0.0-20220927234935-2a059c20f221
