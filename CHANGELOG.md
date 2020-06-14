# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0] (Unreleased) - YYYY-MM-DD

### Added

- GEN: Added a changelog!
- GEN: CI checks that CLI docs are up to date
- GEN: CI does `goreleaser check`
- GEN: Releases for npmfs, docker, scoop
- GEN: Code coverage in code climate
- GEN: Badges & integrations with code climate, codacy, goreportcard, codebeat
- GEN: Compliance with golangci-lint and integration linters
- CLI: Completion commands for bash/fish/zsh/powershell.
- CLI: Full test coverage.
- API: Full test coverage.
- API: New destroy command that matches v2 kv secrets enginve

### Changed

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
