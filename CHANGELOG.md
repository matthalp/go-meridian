# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0] - 2024-10-30

### Added
- Automatic timezone package generation from `timezones.yaml` configuration
- New `timezones/` directory structure for generated timezone packages
- `generate_at_root` boolean option in `timezones.yaml` for backwards compatibility
- Dual-location generation support: existing packages at root and in `timezones/`
- Generator tool at `cmd/generate-timezones/main.go`

### Changed
- **BREAKING**: New timezone packages now generated in `timezones/` directory by default
- Timezone packages are now auto-generated from configuration instead of manually written
- Updated project structure to include `timezones/` directory

### Deprecated
- Root-level timezone packages are maintained for backwards compatibility but new timezones should use `timezones/` directory

## [1.0.0] - 2024-10-14

### Added
- Initial package structure with type-safe timezone handling
- Core `meridian.Time[TZ]` generic type
- Built-in timezone packages: UTC, EST, PST
- `Moment` interface for flexible timezone conversions
- Comprehensive test coverage with race detection
- CI/CD pipeline with GitHub Actions
- Example usage program in `cmd/example`
- Package documentation and examples

## [Unreleased]

### Added
- Nothing yet

### Changed
- Nothing yet

### Deprecated
- Nothing yet

### Removed
- Nothing yet

### Fixed
- Nothing yet

### Security
- Nothing yet
