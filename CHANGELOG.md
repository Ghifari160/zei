<!-- markdownlint-disable MD024 -->
# Changelog

All notable changes in Zei will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [UNRELEASED]

### Added

- Added type `Client`, which implements the same public functions as [`net/http.Client`](https://pkg.go.dev/net/http#Client).
  `Client` can be a drop-in replacement of `net/http.Client`.
- Added `ClientInterface` to assist migration from [`net/http.Client`](https://pkg.go.dev/net/http#Client) to `Client`.

### Changed

### Deprecated

### Removed

### Fixed

### Security
