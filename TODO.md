# TODO

## Soon

- api context and cancellations?
- cli flags for timeouts/context

## Upcoming

- Improve API comments and CLI docs (mention v1/v2 differences and vaku/vault differences)
- Benchmarks
- Example functions
- API example in readme
- doc.go in api
- graceful worker shutdown <https://callistaenterprise.se/blogg/teknik/2019/10/05/go-worker-cancellation/>
- CI to make sure api Version() stays up to date with tags
- CI for line wrapping

## Right Before Release

- www html linting, best practices, css?
- update website
- release on snapcraft
- make sure all 100% on codeclimate, codacy, codebeat, etc...
- make sure changelog is up to date
- review changelog
- update goreleaser.yml with all new checks

## After Release

- Update codebeat, codacy, codelimate to point at main branch
- Rename default branch to main
