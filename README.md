# netatmo-exporter [![Docker Build Status](https://img.shields.io/docker/build/xperimental/netatmo-exporter.svg?style=flat-square)](https://hub.docker.com/r/xperimental/netatmo-exporter/)

Simple [prometheus](https://prometheus.io) exporter for getting sensor values [NetAtmo](https://www.netatmo.com) sensors into prometheus.

## Installation

### Run docker container

The exporter is available as a Docker image: [`xperimental/netatmo-exporter`](https://hub.docker.com/r/xperimental/netatmo-exporter/)

The `latest` tag is built from the current master, tags tagged since the Docker support was added are also available as a tag in Docker.

### Build from source

Because this program uses the "Go Module" feature introduced in Go 1.11, you'll need at least that version of Go for building it.

If you have a working Go installation, getting the binary should be as simple as

```bash
git clone https://github.com/xperimental/netatmo-exporter
cd netatmo-exporter
go build .
```

There is also a `build-arm.sh` script if you want to run the exporter on an ARMv7 device.

## NetAtmo client credentials

This application tries to get data from the NetAtmo API. For that to work you will need to create an application in the [NetAtmo developer console](https://dev.netatmo.com/dev/myaccount), so that you can get a Client ID and secret.

## Usage

```plain
$ netatmo-exporter --help
Usage of netatmo-exporter:
  -a, --addr string            Address to listen on. (default ":9210")
  -i, --client-id string       Client ID for NetAtmo app.
  -s, --client-secret string   Client secret for NetAtmo app.
  -p, --password string        Password of NetAtmo account.
  -u, --username string        Username of NetAtmo account.
```

After starting the server will offer the metrics on the `/metrics` endpoint, which can be used as a target for prometheus.

### Passing secrets

You can pass credentials either via command line arguments (see next section) or by populating the following environment variables:

* `NETATMO_EXPORTER_ADDR` Address to listen on
* `NETATMO_CLIENT_ID` Client ID for NetAtmo app
* `NETATMO_CLIENT_SECRET` Client secret for NetAtmo app
* `NETATMO_CLIENT_USERNAME` Username of NetAtmo account
* `NETATMO_CLIENT_PASSWORD` Password of NetAtmo account

### Scrape interval

The exporter will query the Netatmo API every time it is scraped by prometheus. It does not make sense to scrape the Netatmo API with a small interval as the sensors only update their data every few minutes, so don't forget to set a slower scrape interval for this exporter:

```yml
scrape_configs:
  - job_name: 'netatmo'
    scrape_interval: 90s
    static_configs:
      - targets: ['localhost:9210']
```