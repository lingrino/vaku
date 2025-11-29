# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## 2.11.0 - Upcoming

- API: Add `PathReadMeta` and `PathReadVersion`
- API: Add `PathCopyAllVersions` and `PathMoveAllVersions`
- API: Add `FolderCopyAllVersions` and `FolderMoveAllVersions`
- CLI: Add `--all-versions` flag to path and folder copy/move to include past versions
- CLI: Add `--destroy` flag to path and folder move to permanently delete versions from source
- WWW: Simplify website design

## 2.10.0 - 2025-11-27

### Added

- API: Add `StaticMountProvider` to bypass `sys/mounts` lookup when mount path is known
- CLI: Add `--mount-path` flag to specify mount path directly (bypasses `sys/mounts` permission requirement)
- CLI: Add `--mount-version` flag to specify KV version (1 or 2, defaults to 2)

## 2.9.0 - 2025-09-26

### Changed

- GEN: Vaku is now published as a homebrew cask instead of a formula. You can make the switch by running `brew rm lingrino/tap/vaku && brew install --cask lingrino/tap/vaku`
- GEN: Update to go 1.25
- GEN: Update dependencies

## 2.8.3 - 2025-02-17

### Changed

- GEN: Update to go 1.24
- GEN: Update dependencies

## 2.8.2 - 2024-11-23

### Changed

- GEN: Email tag updates

## 2.8.1 - 2024-11-23

### Changed

- GEN: Update dependencies
- GEN: Update release action

## 2.8.0 - 2024-10-27

### Added

- API: Add ignoreAccessErrors option to skip list and read errors and continue operation
- CLI: Add --ignore-access-errors flag to skip list and read errors and continue operation

### Changed

- GEN: Update dependencies

## 2.7.1 - 2024-03-11

### Changed

- GEN: Update to go 1.23
- GEN: Update dependencies

## 2.7.0 - 2024-03-11

### Changed

- GEN: Remove dependency on Vault import
- GEN: Update dependencies

## 2.6.3 - 2024-02-19

### Changed

- GEN: Update to go 1.22
- GEN: Update dependencies

## 2.6.2 - 2023-08-09

### Changed

- GEN: Update to go 1.21
- GEN: Update dependencies

## 2.6.1 - 2023-02-21

### Changed

- GEN: Update to go 1.20
- GEN: Update dependencies

## 2.6.0 - 2022-11-27

### Changed

- GEN: Update dependencies
- API: Ability to specify custom mount provider. Thank you [@tobgu](https://github.com/tobgu)!

## 2.5.1 - 2022-08-04

### Changed

- GEN: Update to go 1.19

## 2.5.0 - 2022-04-08

### Changed

- GEN: Update dependencies
- CLI: Add `folder write` command

## 2.4.6 - 2022-03-24

### Changed

- GEN: Update to go 1.18
- GEN: Update dependencies

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
