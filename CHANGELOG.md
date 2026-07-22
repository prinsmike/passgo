# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- README badges for the latest release and the module's Go version.

## [2.0.2] - 2026-07-22

### Added
- `CHANGELOG.md`; the release workflow now publishes the matching version's
  changelog section as the GitHub release notes.
- A `verify-changelog` CI job that fails a `v*` tag build when `CHANGELOG.md`
  has no entry for the tag, so releases can't ship auto-generated notes by
  mistake.

## [2.0.1] - 2026-07-22

### Added
- Tag-triggered release workflow that cross-compiles `cmd/passgo` for
  linux, darwin, and windows (amd64/arm64), packages each as an archive
  with a `SHA256SUMS.txt` manifest, and attaches them to the GitHub release.
- README "Download a prebuilt binary" section with a per-platform asset
  table and checksum-verification instructions.

## [2.0.0] - 2026-07-22

Major modernization release. **Breaking changes** — the public API has been
redesigned and the module path is now `github.com/prinsmike/passgo/v2`.

### Security
- All randomness now comes from `crypto/rand` instead of `math/rand`, making
  the output suitable for security-sensitive use.
- Removed the per-character `rand.Seed(time.Now()...)` calls (deprecated since
  Go 1.20 and harmful to unpredictability).

### Added
- `go.mod`; the module path is now `github.com/prinsmike/passgo/v2`.
- A runnable CLI: `go install github.com/prinsmike/passgo/v2/cmd/passgo@latest`.
- Functional-options constructor `New(...Option)` with sensible defaults, so
  `passgo.New()` works out of the box.
- Full test suite (composition, errors, capitalization, concurrency,
  uniqueness) plus a benchmark.
- GitHub Actions CI running gofmt, `go vet`, `golangci-lint` (stricter set),
  and `go test -race`.

### Changed
- `NewPassword` is now `Password(length, numbers, specials)` and performs real
  input validation.
- `Generator` is stateless between calls and safe for concurrent use (the
  shared internal buffer was removed).
- Errors are idiomatic and wrapped.

## [1.0.0] - 2013-03-06

### Added
- Initial release: a human-readable ("pronounceable") password generator based
  on Pradeep Kishore Gowda's `nicepass.py`.

### Security
- This version uses `math/rand` and is **not** cryptographically secure. Use
  v2 or later for any real-world password generation.

[Unreleased]: https://github.com/prinsmike/passgo/compare/v2.0.2...HEAD
[2.0.2]: https://github.com/prinsmike/passgo/compare/v2.0.1...v2.0.2
[2.0.1]: https://github.com/prinsmike/passgo/compare/v2.0.0...v2.0.1
[2.0.0]: https://github.com/prinsmike/passgo/compare/v1.0.0...v2.0.0
[1.0.0]: https://github.com/prinsmike/passgo/releases/tag/v1.0.0
