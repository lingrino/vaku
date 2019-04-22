.PHONY: fmt docs test release

fmt:
	gofmt -l -w -s vaku/

docs:
	go run doc/generate.go

test:
	docker-compose up -d
	go test -cover -race -v ./...

release:
	goreleaser
