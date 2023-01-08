# netatmo-exporter

Simple [prometheus](https://prometheus.io) exporter for getting sensor values [NetAtmo](https://www.netatmo.com) sensors into prometheus.

## Installation

### Run docker container

The exporter is available as a Docker image both on DockerHub and GitHub:

- [`ghcr.io/xperimental/netatmo-exporter`](https://github.com/xperimental/netatmo-exporter/pkgs/container/netatmo-exporter)
- [`xperimental/netatmo-exporter`](https://hub.docker.com/r/xperimental/netatmo-exporter/)

The following tags are available:

- `x.y.z` pointing to the release with that version
- `latest` pointing to the most recent released version
- `master` pointing to the latest build from the default branch

### Build from source

Because this program uses the "Go Module" feature introduced in Go 1.11, you'll need at least that version of Go for building it.

If you have a working Go installation, getting the binary should be as simple as

```bash
git clone https://github.com/xperimental/netatmo-exporter
cd netatmo-exporter
make
```

If you want to build the exporter for a different OS or architecture, you can specify arguments to the Makefile:

```bash
# For 32-bit ARM on Linux
make GO_ARCH=arm
# For 64-bit ARM on Linux
make GO_ARCH=arm64
```

## NetAtmo client credentials

This application tries to get data from the NetAtmo API. For that to work you will need to create an application in the [NetAtmo developer console](https://dev.netatmo.com/apps/), so that you can get a Client ID and secret.

## Usage

```plain
$ netatmo-exporter --help
Usage of netatmo-exporter:
  -a, --addr string                 Address to listen on. (default ":9210")
      --age-stale duration          Data age to consider as stale. Stale data does not create metrics anymore. (default 30m0s)
  -i, --client-id string            Client ID for NetAtmo app.
  -s, --client-secret string        Client secret for NetAtmo app.
      --log-level level             Sets the minimum level output through logging. (default info)
  -p, --password string             Password of NetAtmo account.
      --refresh-interval duration   Time interval used for internal caching of NetAtmo sensor data. (default 8m0s)
  -u, --username string             Username of NetAtmo account.
```

After starting the server will offer the metrics on the `/metrics` endpoint, which can be used as a target for prometheus.

### Passing secrets

You can pass credentials either via command line arguments (see next section) or by populating the following environment variables:

* `NETATMO_EXPORTER_ADDR` Address to listen on
* `NETATMO_CLIENT_ID` Client ID for NetAtmo app
* `NETATMO_CLIENT_SECRET` Client secret for NetAtmo app
* `NETATMO_CLIENT_USERNAME` Username of NetAtmo account
* `NETATMO_CLIENT_PASSWORD` Password of NetAtmo account

### Cached data

The exporter has an in-memory cache for the data retrieved from the Netatmo API. The purpose of this is to decouple making requests to the Netatmo API from the scraping interval as the data from Netatmo does not update nearly as fast as the default scrape interval of Prometheus. Per the Netatmo documentation the sensor data is updated every ten minutes. The default "refresh interval" of the exporter is set a bit below this (8 minutes), but still much higher than the default Prometheus scrape interval (15 seconds).

You can still set a slower scrape interval for this exporter if you like:

```yml
scrape_configs:
  - job_name: 'netatmo'
    scrape_interval: 90s
    static_configs:
      - targets: ['localhost:9210']
```

## Links

- [Grafana Dashboard](https://grafana.com/grafana/dashboards/13672) contributed by [@GordonFreemanK](https://github.com/GordonFreemanK)
