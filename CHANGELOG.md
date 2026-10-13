# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [2.1.0] - 2026-10-13

### Deprecated
- Packages `timezones/est` and `timezones/pst`. They are the same timezones as
  `et` and `pt` (America/New_York and America/Los_Angeles), but their names
  suggest fixed standard-time offsets while they follow daylight saving time.
  Use `et` and `pt` instead. Both will be removed in v3. (#36)

### Fixed
- `Time[TZ].UnmarshalBinary` and `Time[TZ].GobDecode` now store the decoded
  time in UTC. Data encoded from a `time.Time` with a non-UTC offset used to
  leave that offset inside the `Time[TZ]`, so `UTC()` returned a non-UTC time
  and `==` disagreed with an equal value. (#34)
- Timezone package docs now state the offsets each zone uses over the year.
  They previously said every zone observed its name "depending on the time of
  year", which was wrong for zones without daylight saving time and hid that
  `gmt` (Europe/London) is BST, UTC+1, in summer. (#36)
- README and USAGE examples used import paths and functions that no longer
  exist (root-level timezone packages, `Convert`); every full example now
  compiles. Removed merge-conflict markers committed in the README. (#35)

### Changed
- If a timezone database is missing, the panic at package initialization now
  names the package and suggests importing `time/tzdata`. (#36)
- CI runs on every pull request and tests Go 1.20 (the go.mod minimum),
  oldstable and stable; golangci-lint is pinned to v2.14.0. (#33)
- Releases are cut by the Release workflow, which checks the version,
  changelog and API compatibility before creating the tag. (#37)

## [2.0.1] - 2026-06-07

### Fixed
- `Time[TZ].AddDate` now performs calendar arithmetic in the timezone's
  location instead of UTC, preserving the wall-clock time across Daylight
  Saving Time transitions (#29). Previously, adding a day always advanced the
  underlying UTC instant by a literal 24 hours, which shifted the displayed
  wall-clock time by an hour when crossing a DST boundary.

## [2.0.0] - 2025-12-18

### Changed
- **BREAKING**: Module path is now `github.com/matthalp/go-meridian/v2`. (#28)
- **BREAKING**: Timezone packages moved to the `timezones/` directory, e.g.
  `github.com/matthalp/go-meridian/v2/timezones/et`. The root-level packages
  were removed. (#25, #28)

### Added
- `est` and `pst` packages, as aliases of `et` and `pt`, under `timezones/`.
  (Deprecated in 2.1.0.)

## [1.7.0] - 2025-10-20

### Added
- Timezone packages `aest`, `brt`, `cst`, `ct`, `gmt`, `hkt`, `ist`, `jst`,
  `mt` and `sgt`. (#24)

## [1.6.0] - 2025-10-20

### Added
- Timezone package `cet` (Europe/Paris). (#23)

## [1.5.0] - 2025-10-19

### Changed
- **BREAKING** (released as a minor version): Packages `est` and `pst` were
  renamed to `et` and `pt`, because they follow daylight saving time. (#21)

Versions 1.2.0 through 1.4.0 were not released.

## [1.1.0] - 2025-10-14

### Added
- Generic constructors `Parse`, `Unix`, `UnixMilli` and `UnixMicro` in the
  core package, and matching functions in each timezone package. (#18)
- Timezone packages are generated from `timezones.yaml` by
  `cmd/generate-timezones`. (#16)

### Changed
- **BREAKING** (released as a minor version): Timezone packages' `Convert`
  function was renamed to `FromMoment`. (#18)

### Known issues
- The `Version` constant in this release still reads `"1.0.0"`.

## [1.0.0] - 2025-10-13

### Added
- Core `meridian.Time[TZ]` generic type with the `Timezone` and `Moment`
  interfaces.
- Timezone packages `utc`, `est` and `pst`.
- Formatting, arithmetic, comparison and component methods mirroring
  `time.Time`.
- JSON, text, binary and gob encoding, and `database/sql` support
  (`driver.Valuer` and `sql.Scanner`).
- Example program in `cmd/example`.
