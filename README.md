# netatmo-exporter

Simple [prometheus](https://prometheus.io) exporter for getting sensor values [NetAtmo](https://www.netatmo.com) sensors into prometheus.

## Installation

If you have a working Go installation, getting the binary should be as simple as

```bash
go get github.com/xperimental/netatmo-exporter
```

There is also a `build-arm.sh` script if you want to run the exporter on an ARMv7 device.

## NetAtmo client credentials

This application tries to get data from the NetAtmo API. For that to work you will need to create an application in the [NetAtmo developer console](https://dev.netatmo.com/dev/myaccount), so that you can get a Client ID and secret.

## Usage

```
$ netatmo-exporter --help
Usage of netatmo-exporter:
  -a, --addr string            Address to listen on. (default ":8080")
  -i, --client-id string       Client ID for NetAtmo app.
  -s, --client-secret string   Client secret for NetAtmo app.
  -p, --password string        Password of NetAtmo account.
  -u, --username string        Username of NetAtmo account.
```

After starting the server will offer the metrics on the `/metrics` endpoint, which can be used as a target for prometheus.

The exporter will query the Netatmo API every time it is scraped by prometheus. It does not make sense to scrape the Netatmo API with a small interval as the sensors only update their data every few minutes, so don't forget to set a slower scrape interval for this exporter:

```yml
scrape_configs:
  - job_name: 'netatmo'
    scrape_interval: 90s
    static_configs:
      - targets: ['localhost:8080']
```

**Note:** The exporter currently uses port 8080 as a default as it does not have an "assigned exporter port" yet. Look at the [prometheus Wiki](https://github.com/prometheus/prometheus/wiki/Default-port-allocations) for any updates.
