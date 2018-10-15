# vaku

[![Vaku](www/assets/logo-vaku-sm.png?raw=true)](www/assets/logo-vaku-sm.png "Vaku")

[![CircleCI](https://circleci.com/gh/Lingrino/vaku.svg?style=svg)](https://circleci.com/gh/Lingrino/vaku)
[![Go Report Card](https://goreportcard.com/badge/github.com/Lingrino/vaku)](https://goreportcard.com/report/github.com/Lingrino/vaku)
[![GoDoc](https://godoc.org/github.com/Lingrino/vaku/vaku?status.svg)](https://godoc.org/github.com/Lingrino/vaku/vaku)

A CLI and Go API that add useful functions on top of Hashicorp Vault.

Please read the [godoc documentation](https://godoc.org/github.com/Lingrino/vaku/vaku)
for all API usage information and examples.

Please use `vaku help` in your terminal for all documentation and usage information
regarding the Vaku CLI

Vaku is now V1. The API and CLI will be backwards compatible until the next point release.
See the checklist below for progress and upcoming features.

**API/CLI Functionality:**

- [x] Path List
- [x] Path Read
- [x] Path Write (API only)
- [x] Path Delete
- [x] Path Destroy (v2 mounts only)
- [x] Path Copy
- [x] Path Move
- [x] Path Update (API only)
- [x] Path Search
- [ ] Path Diff
- [x] Folder List
- [x] Folder Read
- [x] Folder Write (API only)
- [x] Folder Delete
- [x] Folder Destroy (v2 mounts only)
- [x] Folder Copy
- [x] Folder Move
- [x] Folder Search
- [ ] Folder Diff
- [x] Folder Map (CLI Only)
- [ ] Add Timeouts to Workers

**CLI Improvements:**

- [ ] Add tests
- [ ] Add to homebrew
- [ ] Add example usage
- [ ] Support concurrency flag

**Running Tests:**

Tests are meant to be run side by side with a real Vault server docker image. This
creates an external dependency for the tests but makes it much simpler to test different
Vault versions and key/value mounts. With docker and docker-compose installed tests
can be run with a simple `make test`. CircleCI will also build all commits and report
status on all PRs.
