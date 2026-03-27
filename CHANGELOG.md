# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachanglog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.4.2] - 2026-03-28

### Changed
- Unified Makefile: replaced macOS universal binary with separate `darwin/amd64` and `darwin/arm64` targets; standardized targets (`build`, `build-all`, `test`, `lint`, `check`, `package`, `clean`, `help`) and output layout (`dist/` flat directory, `.zip` archives).

## [1.4.1] - 2026-03-28

### Internal

- Updated Go module path to `github.com/nlink-jp/json-to-table` following repository transfer to nlink-jp organization.

## [1.4.0] - 2025-08-28

### Added

- Added CSV output format (`--format csv`).

## [1.3.2] - 2025-08-26

### Fixed

- Fixed an issue where full-width alphabets were garbled in PNG output by replacing the font with `MPLUS1p-Regular.ttf`.

## [1.3.1] - 2025-08-14

### Changed

- Refactored `json-to-table.go` by splitting it into multiple, smaller files for improved modularity and maintainability.
  - `main.go`: Application entry point.
  - `core.go`: Core parsing logic (`parseJSON`, `matchHeaders`).
  - `render_text.go`: Text and Markdown rendering.
  - `render_image.go`: PNG image rendering.
  - `render_html.go`: HTML rendering.
  - `render_slack.go`: Slack Block Kit rendering.

## [1.3.0] - 2025-08-14

### Added

- Reimplemented column exclusion feature with `--exclude-columns` (`-e`) flag.
  - Exclusion is processed before column inclusion (`--columns`).
  - Supports specific column names and wildcard patterns (`*`, `prefix*`).

### Changed

- Refactored column selection logic in `parseJSON` to support exclusion precedence.
- Moved `test_data.json` to `testdata/test_data.json` for better organization.
- Updated `README.md` and `README.ja.md` to reflect new column selection flags and `testdata` usage.

## [1.2.0] - 2025-08-14

### Added

- Added 'blocks' as a shorthand for 'slack-block-kit' output format (`--format blocks`).

## [1.1.0] - 2025-08-14

### Added

- Slack Block Kit output format (`--format slack-block-kit`).

### Fixed

- Removed redundant 'v' prefix from package filenames.

## [1.0.0] - 2025-08-05

### Added

- HTML output format (`--format html`).
- Version information (`--version`) via `ldflags`.
- `FONTS_LICENSE` file for Mplus 1 Code font.
- `README.ja.md` for Japanese documentation.
- `make package` target to create zipped release archives.

### Changed

- Build system switched from `build.sh` to `Makefile`.
- Build output directory changed from `dist_table` to `dist`.
- Standardized executable names (e.g., `json-to-table`, `json-to-table.exe`).
- Font handling: Mplus 1 Code font is now embedded directly in the repository.
- READMEs updated to reflect new build process and HTML output.
- Overview in READMEs clarified to state `splunk-cli` companion tool role.
- Improved error handling with `%+v` formatting for better debuggability.

### Removed

- Outdated `BUILD.md` file.
- `build.sh` script.

### Fixed

- Corrected `parseJSON` return value for empty data.
- Resolved font download issues during build by embedding the font directly.
- Fixed redundant 'v' prefix in release package filenames.
