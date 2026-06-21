# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [3.0.0] - 2026-06-20

### Added
- `Time[TZ].To[TZ2]()` generic method for fluent, type-safe timezone
  conversion: `eastern.To[pt.Timezone]()` is the method form of
  `pt.FromMoment(eastern)` and chains naturally
  (`t.To[utc.Timezone]().To[pt.Timezone]()`). This relies on generic methods,
  a language feature introduced in Go 1.27.

### Changed
- **BREAKING**: Module path is now `github.com/matthalp/go-meridian/v3`.
- **BREAKING**: Minimum supported Go version is now 1.27 (was 1.20), required
  by the generic `To` method.
- Migrated the golangci-lint configuration to the v2 schema and updated CI to
  build and test against Go 1.27.

### Notes
- `FromMoment` is unchanged and remains the way to convert from a plain
  `time.Time` (which has no `To` method). The `To` method is additive sugar for
  `meridian.Time[TZ]` → `meridian.Time[TZ2]` conversions.

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

## [2.0.1] - 2026-06-07

### Fixed
- `Time[TZ].AddDate` now performs calendar arithmetic in the timezone's
  location instead of UTC, preserving the wall-clock time across Daylight
  Saving Time transitions (#29). Previously, adding a day always advanced the
  underlying UTC instant by a literal 24 hours, which shifted the displayed
  wall-clock time by an hour when crossing a DST boundary.
