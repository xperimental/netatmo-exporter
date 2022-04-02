# Changelog

This changelog contains the changes made between releases. The versioning follows [Semantic Versioning](https://semver.org/).

## Unreleased

### Changed

- Go 1.17

### Fixed

- lastRefreshError is not reset (#11)
- Docker build for arm64

## 1.3.0 (2020-08-09)

### Added

- HTTP Handler for getting build information `/version`
- In-memory cache for data retrieved from NetAtmo API, configurable timeouts

### Changed

- Logger uses leveled logging, added option to set log level
- Updated Go runtime and dependencies

## 1.2.0 (2018-10-27)

### Added

- Support for battery and RF-link status
- Support for configuration via environment variables

## 1.1.0 (2018-09-02)

### Added

- Support for wind and rain sensors

### Changed

- Metrics now also contain a label for the "station name"

## 1.0.1 (2017-11-26)

### Fixed

- Integrate fix of upstream library

## 1.0.0 (2017-03-09)

- Initial release
