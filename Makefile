.PHONY: fmt docs test release install

fmt:
	gofmt -l -w -s vaku/

docs:
	go run doc/generate.go

test:
	docker-compose up -d
	go test -cover -v ./...

release:
	goreleaser

install:
	go install github.com/lingrino/vaku
