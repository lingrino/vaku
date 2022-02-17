# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## 2.4.5 - 2022-02-17

Thank you [@szechuen](https://github.com/szechuen) for fixing two important worker issues!

### Changed

- GEN: Update dependencies
- API: Fix folder list errgroup cancellation
- API: Fix workers preamturely returning on success

## 2.4.4 - 2022-01-28

### Changed

- GEN: Update dependencies

## 2.4.3 - 2022-01-06

### Changed

- GEN: Update dependencies

## 2.4.2 - 2021-11-17

### Changed

- GEN: Update to go 1.17
- GEN: Update dependencies

## 2.4.1 - 2021-11-08

### Changed

- GEN: Package upgrades.
- GEN: New release fixes homebrew warnings.

## 2.4.0 - 2021-07-08

### Added

- CLI: Use built in completion commands from cobra 1.2.0
- GEN: Package updates. A fresh `go get` now completes without error.

## 2.3.0 - 2021-05-15

Thank you to [@karakanb](https://github.com/karakanb) for finding and fixing a tricky bug!

### Changed

- API: Add `AddPrefix` and `AddPrefixList` helper functions.
- GEN: Fixed a bug where `folder search` would hang if the mount path shared a name with a folder.

## 2.2.1 - 2021-05-15

### Changed

- GEN: Update dependencies
- GEN: Run actions on PRs to support forked contributions
- GEN: Update golangci-lint version and configuration
- API: Fix new golangci-lint error wrapping issues

## 2.2.0 - 2021-02-16

### Changed

- GEN: Update to go 1.16
- GEN: Support arm macs
- GEN: Update dependencies

## 2.1.2 - 2020-11-11

### Changed

- GEN: Upgrade all go dependencies

## 2.1.1 - 2020-10-31

### Added

- GEN: Update packages
- GEN: Update golangci-lint
- GEN: Enable `errorlint`, `tparallel`, `wrapcheck` linters
- GEN: Fix found linter issues
- API: `ErrApplyOptions` now returned in `api.NewClient`

## 2.1.0 - 2020-10-27

### Added

- CLI: Add flags and support for vault namespaces

## 2.0.0 - 2020-10-17

### Added

- GEN: Added a changelog!
- GEN: CI checks that CLI docs are up to date
- GEN: CI enforces `go mod tidy`
- GEN: CI checks for vaku.Version() matching tagged version
- GEN: CI does `goreleaser check`
- GEN: Releases for npmfs, docker, scoop
- GEN: Code coverage in code climate
- GEN: Badges & integrations with code climate, codacy, goreportcard, codebeat
- GEN: Compliance with golangci-lint and integration linters
- GEN: The default branch name is now `main`
- GEN: Examples in readme
- API: Full test coverage.
- API: New destroy command that matches v2 kv secrets engine
- CLI: Full test coverage.
- CLI: Completion commands for bash/fish/zsh/powershell.
- CLI: Flag for sorting output
- CLI: Flag support for separate source/dest vault servers
- CLI: Hidden commands for unsupport CLI calls will redirect users to the API

### Changed

- API: The api package is now `vaku/api` instead of `vaku/vaku`
- API: Concurrency limits now set on the client.
- API: All errors are now exported and can be unwrapped.
- API: Tests now use an inmem vault server instead of docker.
- API: Tests can be run directly with `go test`.
- API: Client is now configured using functional options.
- API: Client now supports source/destination vault clients.
- API: CLient no longer inherits the vault client. Set in source/dest instead.
- API: Destory calls renamed to DeleteMeta which is more accurate.

### Removed

- API: Removed `PathInput{}`. Functions now take a path string
- API: Removed public mount functions
- API: Removed unused public helper functions

## 1.x.x

The changelog was started with the `2.0.0` release, which was a complete rewrite. A record of `1.x.x` changes can be found in [GitHub releases](https://github.com/lingrino/vaku/releases) and/or git commit history.
