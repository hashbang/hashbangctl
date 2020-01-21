NAMESPACE ?= hashbangctl

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
test: \
	docker-build-test \
	docker-restart \
	docker-test \
	docker-stop

.PHONY: test-shell
test-shell: \
	docker-build-test \
	docker-restart \
	docker-test-shell \
	docker-stop

.PHONY: clean
clean: docker-clean
	rm -rf ./go ./bin

.PHONY: fetch
fetch:
	git submodule update --init --recursive

.PHONY: fetch-latest
fetch-latest:
	git submodule foreach 'git checkout master && git pull'

## Secondary Targets

.PHONY: docker-build
docker-build:
	docker build -t local/$(NAMESPACE) .
	docker build -t local/$(NAMESPACE)-userdb test/modules/userdb/
	docker build \
		--build-arg=POSTGREST_VERSION=v6.0.2 \
		-t local/$(NAMESPACE)-postgrest \
		test/modules/postgrest/docker/

.PHONY: docker-restart
docker-restart: docker-stop docker-start

.PHONY: docker-start
docker-start:
	docker network inspect $(NAMESPACE) \
	|| docker network create $(NAMESPACE)
	docker inspect -f '{{.State.Running}}' $(NAMESPACE) 2>/dev/null \
	|| docker run \
		--detach=true \
		--name $(NAMESPACE) \
		--network=$(NAMESPACE) \
		--expose="2222" \
		-p "2222:2222" \
		local/$(NAMESPACE)
	docker inspect -f '{{.State.Running}}' $(NAMESPACE)-userdb 2>/dev/null \
	|| docker run \
		--rm \
		--detach=true \
		--name $(NAMESPACE)-userdb \
		--network=$(NAMESPACE) \
		--env="CONTAINER=$(NAMESPACE)" \
		local/$(NAMESPACE)-userdb
	docker inspect -f '{{.State.Running}}' $(NAMESPACE)-postgrest 2>/dev/null \
	|| docker run \
		--rm \
		--detach=true \
		--name $(NAMESPACE)-postgrest \
		--network=$(NAMESPACE) \
		--env="CONTAINER=$(NAMESPACE)" \
		local/$(NAMESPACE)-postgrest

.PHONY: docker-stop
docker-stop:
	docker inspect -f '{{.State.Running}}' $(NAMESPACE) 2>/dev/null \
	&& docker rm -f $(NAMESPACE) || true
	docker inspect -f '{{.State.Running}}' $(NAMESPACE)-userdb 2>/dev/null \
	&& docker rm -f $(NAMESPACE)-userdb || true
	docker inspect -f '{{.State.Running}}' $(NAMESPACE)-postgrest 2>/dev/null \
	&& docker rm -f $(NAMESPACE)-postgrest || true

.PHONY: docker-log
docker-log:
	docker logs -f $(NAMESPACE)

.PHONY: docker-clean
docker-clean: docker-stop
	docker image rm local/$(NAMESPACE)
	docker image rm local/$(NAMESPACE)-test
	docker image rm local/$(NAMESPACE)-postgrest
	docker image rm local/$(NAMESPACE)-userdb

.PHONY: docker-test
docker-test:
	docker run \
		--rm \
		--hostname=$(NAMESPACE)-test \
		--name $(NAMESPACE)-test \
		--network=$(NAMESPACE) \
		--env="CONTAINER=$(NAMESPACE)" \
		local/$(NAMESPACE)-test

.PHONY: docker-test-shell
docker-test-shell: \
	docker-build docker-stop docker-start docker-build-test docker-stop
	docker run \
		--rm \
		-it \
		--hostname=$(NAMESPACE)-test \
		--name $(NAMESPACE)-test \
		--network=$(NAMESPACE) \
		--env CONTAINER="$(NAMESPACE)" \
		local/$(NAMESPACE)-test \
		bash

.PHONY: docker-build-test
docker-build-test:
	docker build -t local/$(NAMESPACE)-test test/
