# vaku
[![CircleCI](https://circleci.com/gh/Lingrino/vaku.svg?style=svg)](https://circleci.com/gh/Lingrino/vaku)

[![Go Report Card](https://goreportcard.com/badge/github.com/Lingrino/vaku)](https://goreportcard.com/report/github.com/Lingrino/vaku)

Useful functions in Go for Hashicorp Vault.

Please read the [godoc documentation](https://godoc.org/github.com/Lingrino/vaku/vaku)
for all API usage information and examples.

Please use `vaku help` in your terminal for all documentation and usage information
regarding the Vaku CLI

Vaku is now V1. The API and CLI will be backwards compatible until the next point release.
See the checklist below for progress and upcoming features.

**API Functionality:**
- [x] Path List
- [x] Path Read
- [x] Path Write
- [x] Path Delete
- [ ] Path Destroy (v2 mounts only)
- [x] Path Copy
- [x] Path Move
- [x] Path Update
- [x] Path Search
- [ ] Path Diff
- [x] Folder List
- [x] Folder Read
- [x] Folder Write
- [x] Folder Delete
- [ ] Folder Destroy (v2 mounts only)
- [x] Folder Copy
- [x] Folder Move
- [x] Folder Search
- [ ] Folder Diff
- [ ] Folder Map
- [ ] Policy Enforce
- [ ] Approle Enforce
- [ ] Userpass Enforce
- [ ] Add Timeouts to Workers
- [ ] Support Wrapped Secrets

**CLI Improvements:**
- [ ] Add tests
- [ ] Add to homebrew
- [ ] Add example usage
- [ ] Support concurrency flag
- [ ] Support more than JSON output
- [ ] Add write/update commands (native CLI probably better for writing data)

**Running Tests:**

Tests are meant to be run side by side with a real Vault server docker image. This
creates an external dependency for the tests but makes it much simpler to test different
Vault versions and key/value mounts. With docker and docker-compose installed tests
can be run with a simple `make test`. CircleCI will also build all commits and report
status on all PRs.
