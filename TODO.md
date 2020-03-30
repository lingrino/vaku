# TODO

## Soon

- CLI tests
- man pages
- include man pages with cli
- ci to make sure cobra docs are up to date
- bash completion
- add context and cancellations?
- add installation instructions like <https://goreleaser.com/install/>
- add a changelog <https://keepachangelog.com/en/1.0.0/>
  - <https://github.com/starship/starship/releases>
- ci should ignore generated markdown

## Upcoming

- Benchmarks
- Example functions
- API example in readme
- doc.go in api
- further parallelize tests?
- cli option to sort output
- graceful worker shutdown <https://callistaenterprise.se/blogg/teknik/2019/10/05/go-worker-cancellation/>
- cli checks for updates <https://github.com/tcnksm/go-latest>
- CI to make sure api Version() stays up to date with tags

## Right Before Release

- www html linting, best practices, css?
- update website
- release on snapcraft
- verify delete/destory behavior on real vault
- codacy file and style fixes
- make sure all 100% on codeclimate, codacy, codebeat, etc...

## After Release

- Update codebeat, codacy, codelimate to point at master branch
