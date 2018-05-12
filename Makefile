.PHONY: test

test:
	# TODO - dont't kill/sleep if already running
	docker stop test-vault || true && docker rm test-vault || true && \
	docker run -d --name=test-vault -p 8300:8300 \
	-e VAULT_DEV_ROOT_TOKEN_ID=hunter2 \
	-e VAULT_DEV_LISTEN_ADDRESS=0.0.0.0:8300 \
	-e VAULT_LOG=deug \
	vault:latest && sleep 8
	export VAULT_ADDR=http://localhost:8300 && \
	export VAULT_TOKEN=hunter2 && \
	go test -v ./...
