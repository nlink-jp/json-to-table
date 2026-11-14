# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachanglog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.5.1] - 2026-11-14

### Security

- **`golang.org/x/image` is updated from v0.29.0 to v0.46.0.** `govulncheck`
  reported GO-2026-4962 (a crafted font makes the SFNT parser allocate too
  much memory) as reachable, and it is: the PNG output parses a font. The
  only font it parses is the one built into the binary, so no input of the
  user's reaches that parser. Nine more advisories in the same module are in
  image decoders (TIFF, WebP, BMP) this tool does not use. The update takes
  the affected code out of the build.
- **The binaries are built with Go 1.27.4, and `go.mod` now requires it** (it
  was 1.24.5). The binaries of v1.5.0 were built with Go 1.26.4 in July.

### Changed

- **`make package` runs the vulnerability check first** (`make vuln`), so a
  release cannot be built past a finding.
- `make check` passes again: five findings of the current linter (two
  `Close` results not looked at, three `WriteString(fmt.Sprintf(...))`) are
  resolved with no change in behaviour.

### Fixed

- **`make verify-release` now fails closed.** Its last block chained unzip, the
  packaged binary's `--version` and `spctl` with `&&` and ended the whole chain
  in `|| true`, so a zip that did not unpack or a binary that did not run exited
  0 and the upload proceeded. Each step is now judged on its own, the packaged
  binary's `--version` must contain the tag being released, and only the
  informational `spctl` line may be ignored. Matches the org template
  (CONVENTIONS.md §Code Signing → Verifying a release).
- **The Linux archives no longer carry macOS file metadata.** macOS `tar` wrote
  each bundled file's extended attributes (`com.apple.provenance`, and a Dropbox
  attribute where the tree is synced) into the `.tar.gz` twice: as AppleDouble
  `._` members, which GNU tar extracts as stray `._<name>` files beside the real
  ones, and as `LIBARCHIVE.xattr.*` / `SCHILY.xattr.*` pax headers, which it
  reports as unknown keywords. `make package` now archives with
  `COPYFILE_DISABLE=1 tar --no-xattrs`; each setting stops one of the two.
  Archives already published still carry them; the files themselves are
  unaffected.

### Internal

- `make verify-release` also judges each Linux archive: no AppleDouble or other
  macOS metadata members — listed with `--options 'tar:!mac-ext'`, because a
  plain macOS listing folds `._` members away — no extended attributes as pax
  headers, and exactly the canonical binary, `README.md`, `LICENSE` and
  `FONTS_LICENSE`, compared in the C locale.
- The Linux-archive check in `make verify-release` reads each archive's pax
  headers with Python's `tarfile` instead of grepping the decompressed stream,
  which also matched file text that names the keywords (a bundled CHANGELOG,
  for one).

## [1.5.0] - 2026-07-12

### Removed

- **darwin/amd64 (Intel) pre-built binary.** macOS releases now ship
  **arm64 only**, per the org-wide policy (darwin is Apple-Silicon only; no
  universal binaries). Intel Mac users can build from source.

### Changed

- **Linux release archives are now `.tar.gz`** (darwin/windows remain `.zip`),
  per `nlink-jp/.github` CONVENTIONS.md §Release Archive Standard. Archives
  still bundle `LICENSE` + `README.md` + `FONTS_LICENSE` alongside the
  canonical binary.
- **darwin code-signature identifier** is now the canonical `json-to-table`
  (was `json-to-table-darwin-arm64`), set via `codesign -i` so it stays
  stable after the archived binary is renamed to its canonical name.

No change to the binary's behaviour — a packaging / build-config release.

## [1.4.4] - 2026-05-23

### Changed

- **Darwin releases are now Developer ID signed and Apple-notarized.**
  `json-to-table-v1.4.4-darwin-{amd64,arm64}.zip` carry full Apple
  Developer ID Application signatures and notarization tickets from
  Apple. End users on macOS no longer need to bypass Gatekeeper
  with right-click → Open or `xattr -d com.apple.quarantine` on
  first launch; local users who place `json-to-table` under
  Dropbox-synced (or any other FileProvider-managed) paths are no
  longer killed by macOS's ad-hoc + provenance distrust policy.
  Pipeline: `scripts/codesign-darwin.sh` +
  `scripts/notarize-darwin.sh`, driven by `make package`. Adopts
  the org-wide convention in `nlink-jp/.github` CONVENTIONS.md
  §Code Signing.

No behaviour change to the binary itself — feature-wise this is
identical to v1.4.3.

## [1.4.3] - 2026-03-28

### Internal

- Extracted `measureWidth()` helper; eliminated duplicated width-calculation loops in `render_text.go`.
- Extracted `newRichTextCell()` helper; moved Slack type definitions to package scope in `render_slack.go`.
- Replaced pixel-by-pixel grid drawing with `draw.Draw` in `render_image.go`.
- Fixed silently ignored `os.Stdin.Stat()` error in `main.go`.
- Added 31 unit tests covering `measureWidth`, `matchHeaders`, `parseJSON`, and all renderers.

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
