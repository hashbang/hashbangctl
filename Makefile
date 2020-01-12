## Primary Targets

.PHONY: build
build:
	GOBIN=$(PWD)/bin \
	GOPATH=$(PWD)/go \
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64 \
	go install ./...

serve:
	bin/server

.PHONY: connect
connect:
	SSH_AUTH_SOCK="" \
	ssh \
		-i test/keys/id_ed25519 \
		-o UserKnownHostsFile=/dev/null \
		-o StrictHostKeyChecking=no \
		-p2222 localhost

.PHONY: test
test: docker-build docker-start docker-test docker-stop

.PHONY: test-shell
test-shell: docker-start docker-test-shell docker-stop

.PHONY: clean
clean: docker-clean
	rm -rf ./go ./bin

## Secondary Targets

.PHONY: docker-build
docker-build:
	docker build -t local/hashbangctl .

.PHONY: docker-start
docker-start:
	docker network inspect hashbangctl \
	|| docker network create hashbangctl
	docker inspect -f '{{.State.Running}}' hashbangctl 2>/dev/null \
	|| docker run \
		--detach=true \
		--name hashbangctl \
		--network=hashbangctl \
		--expose="2222" \
		-p "2222:2222" \
		local/hashbangctl

.PHONY: docker-stop
docker-stop:
	docker inspect -f '{{.State.Running}}' hashbangctl 2>/dev/null \
	&& docker rm -f hashbangctl || true

.PHONY: docker-log
docker-log:
	docker logs -f hashbangctl

.PHONY: docker-clean
docker-clean: docker-stop
	docker image rm local/hashbangctl

.PHONY: test
docker-test: docker-stop docker-start docker-build-test
	docker run \
		--rm \
		--hostname=hashbangctl-test \
		--name hashbangctl-test \
		--network=hashbangctl \
		--env="CONTAINER=hashbangctl" \
		local/hashbangctl-test

.PHONY: docker-test-shell
docker-test-shell: docker-stop docker-start docker-build-test
	docker run \
		--rm \
		-it \
		--hostname=hashbangctl-test \
		--name hashbangctl-test \
		--network=hashbangctl \
		--env CONTAINER="hashbangctl" \
		local/hashbangctl-test \
		bash

.PHONY: docker-build-test
docker-build-test:
	docker build -t local/hashbangctl-test test/
