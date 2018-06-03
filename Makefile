.PHONY: test

fmt:
	gofmt -l -w -s vaku/

test:
	docker-compose up -d
	go test -cover -v ./...

release:
	goreleaser
