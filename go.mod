module github.com/xperimental/netatmo-exporter/v2

go 1.23

toolchain go1.23.2

require (
	github.com/exzz/netatmo-api-go v0.0.0-20201009073308-a8620474d1ea
	github.com/google/go-cmp v0.6.0
	github.com/prometheus/client_golang v1.20.5
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/pflag v1.0.5
	golang.org/x/oauth2 v0.23.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.60.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	golang.org/x/sys v0.26.0 // indirect
	google.golang.org/protobuf v1.35.1 // indirect
)

replace github.com/exzz/netatmo-api-go => github.com/xperimental/netatmo-api-go v0.0.0-20241020165001-b941db716664
