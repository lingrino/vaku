# Vaku

[![Vaku](www/assets/images/logo-vaku-sm.png?raw=true)](www/assets/logo-vaku-sm.png "Vaku")

[![PkgGoDev](https://pkg.go.dev/badge/github.com/lingrino/vaku/v2/api)](https://pkg.go.dev/github.com/lingrino/vaku/v2/api)
[![goreportcard](https://goreportcard.com/badge/github.com/lingrino/vaku)](https://goreportcard.com/report/github.com/lingrino/vaku)
[![Code Quality](https://app.codacy.com/project/badge/Grade/65802905eb8148e2ae9ae4c909673ee2)](https://www.codacy.com/gh/lingrino/vaku/dashboard)
[![Test Coverage](https://api.codeclimate.com/v1/badges/db6951b0aa53becf8c92/test_coverage)](https://codeclimate.com/github/lingrino/vaku/test_coverage)
[![Maintainability](https://api.codeclimate.com/v1/badges/db6951b0aa53becf8c92/maintainability)](https://codeclimate.com/github/lingrino/vaku/maintainability)

Vaku is a CLI and API for running path- and folder-based operations on the Vault Key/Value secrets engine. Vaku extends the existing Vault CLI and API by allowing you to run the same path-based list/read/write/delete functions on folders as well. Vaku also lets you search, copy, and move both secrets and folders.

## Installation

### Homebrew

```shell
brew install lingrino/tap/vaku
```

### Scoop

```shell
scoop bucket add vaku https://github.com/lingrino/scoop-vaku.git
scoop install vaku
```

### Docker

```shell
docker run ghcr.io/lingrino/vaku --help
```

### Binary

Download the latest binary or deb/rpm for your os/arch from the [releases page](https://github.com/lingrino/vaku/releases).

## Usage

Vaku CLI documentation can be found on the command line using either `vaku help [cmd]` or `vaku [cmd] --help`. The same documentation is also available in markdown form in the [docs/cli](docs/cli/vaku.md) folder.

## API

Documentation for the Vaku API is on [pkg.go.dev](https://pkg.go.dev/github.com/lingrino/vaku/v2/api).

## Contributing

Suggestions and contributions of all kinds are welcome! If there is functionality you would like to see in Vaku please open an Issue or Pull Request and I will be sure to address it.

## Tests

Vaku is well tested and uses only the standard go testing tools.

```shell
$ go test -cover -race ./...
ok  github.com/lingrino/vaku/v2      0.095s coverage: 100.0% of statements
ok  github.com/lingrino/vaku/v2/api 12.065s coverage: 100.0% of statements
ok  github.com/lingrino/vaku/v2/cmd  0.168s coverage: 100.0% of statements
```
