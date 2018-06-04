.PHONY: fmt test release install

fmt:
	gofmt -l -w -s vaku/

test:
	docker-compose up -d
	go test -cover -v ./...

release:
	goreleaser

install:
	go install github.com/Lingrino/vaku
