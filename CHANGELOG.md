# Changelog

This changelog contains the changes made between releases. The versioning follows [Semantic Versioning](https://semver.org/).

## Unreleased

## [2.1.1] - 2024-12-22

### Added

- Error messages returned from NetAtmo API are shown in more detail

### Changed

- Updated Go runtime and dependencies

## [2.1.0] - 2024-10-20

### Added

- Show version information on startup
- Token is saved during runtime of exporter once it is refreshed and not just on shutdown

### Fixed

- Ignore expired tokens on startup

### Changed

- Updated Go runtime and dependencies

## [2.0.1] - 2023-10-15

### Changed

- Maintenance release, updates Go runtime and dependencies

## [2.0.0] - 2023-07-18

- Major: New authentication method replaces existing username/password authentication

## [1.5.1] - 2023-01-08

### Changed

- `latest` Docker tag now points to most recent release and `master` points to the build from the default branch

## [1.5.0] - 2022-12-06

### Added

- Debugging endpoint for looking at data read from NetAtmo API (`/debug/data`)
- New `home` label as additional identification for sensors
- Use module ID (currently MAC-address) as fallback for the `name` label if no name is provided

### Changed

- Switch to fork of netatmo-api-go library

### Fixed

- Not all metric descriptors sent to registry

## [1.4.0] - 2022-04-02

### Changed

- Go 1.17

### Fixed

- Updated Prometheus client library for CVE-2022-21698
- lastRefreshError is not reset (#11)
- Docker build for arm64

## [1.3.0] - 2020-08-09

### Added

- HTTP Handler for getting build information `/version`
- In-memory cache for data retrieved from NetAtmo API, configurable timeouts

### Changed

- Logger uses leveled logging, added option to set log level
- Updated Go runtime and dependencies

## [1.2.0] - 2018-10-27

### Added

- Support for battery and RF-link status
- Support for configuration via environment variables

## [1.1.0] - 2018-09-02

### Added

- Support for wind and rain sensors

### Changed

- Metrics now also contain a label for the "station name"

## [1.0.1] - 2017-11-26

### Fixed

- Integrate fix of upstream library

## [1.0.0] - 2017-03-09

- Initial release

[2.1.1]: https://github.com/xperimental/netatmo-exporter/releases/tag/v2.1.1
[2.1.0]: https://github.com/xperimental/netatmo-exporter/releases/tag/v2.1.0
[2.0.1]: https://github.com/xperimental/netatmo-exporter/releases/tag/v2.0.1
[2.0.0]: https://github.com/xperimental/netatmo-exporter/releases/tag/v2.0.0
[1.5.1]: https://github.com/xperimental/netatmo-exporter/releases/tag/v1.5.1
[1.5.0]: https://github.com/xperimental/netatmo-exporter/releases/tag/v1.5.0
[1.4.0]: https://github.com/xperimental/netatmo-exporter/releases/tag/v1.4.0
[1.3.0]: https://github.com/xperimental/netatmo-exporter/releases/tag/v1.3.0
[1.2.0]: https://github.com/xperimental/netatmo-exporter/releases/tag/v1.2.0
[1.1.0]: https://github.com/xperimental/netatmo-exporter/releases/tag/v1.1.0
[1.0.1]: https://github.com/xperimental/netatmo-exporter/releases/tag/v1.0.1
[1.0.0]: https://github.com/xperimental/netatmo-exporter/releases/tag/v1.0.0
