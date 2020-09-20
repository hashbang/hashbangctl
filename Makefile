NAMESPACE ?= hashbangctl

## Primary Targets

.PHONY: build
build: docker-build

.PHONY: build-native
build-native:
	GOBIN=$(PWD)/bin \
	GOPATH=$(PWD)/go \
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64 \
	go install ./...

.PHONY: serve
serve: docker-start docker-logs docker-stop

.PHONY: serve-native
serve-native:
	API_URL="https://userdb.hashbang.sh/v1" \
	API_TOKEN="eyJhbGciOiJIUzI1NiJ9.eyJyb2xlIjoiYXBpLXVzZXItY3JlYXRlIn0.iOcRzRAjPsT9DOhu5OSeRuQ38D3KL5NppsfyuZYiDeI" \
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
	docker build -t local/$(NAMESPACE)-userdb modules/userdb/
	docker build \
		--build-arg=POSTGREST_VERSION=v6.0.2 \
		-t local/$(NAMESPACE)-postgrest \
		modules/postgrest/docker/

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
		--env API_URL="http://hashbangctl-postgrest:3000" \
		--env API_TOKEN="eyJhbGciOiJIUzI1NiJ9.eyJyb2xlIjoiYXBpLXVzZXItY3JlYXRlIn0.iOcRzRAjPsT9DOhu5OSeRuQ38D3KL5NppsfyuZYiDeI" \
		--expose="2222" \
		-p "2222:2222" \
		local/$(NAMESPACE)
	docker inspect -f '{{.State.Running}}' $(NAMESPACE)-userdb 2>/dev/null \
	|| docker run \
		--rm \
		--detach=true \
		--name $(NAMESPACE)-userdb \
		--network=$(NAMESPACE) \
		--volume $(PWD)/test/test_data.sql:/docker-entrypoint-initdb.d/99-init.sql \
		local/$(NAMESPACE)-userdb
	docker inspect -f '{{.State.Running}}' $(NAMESPACE)-postgrest 2>/dev/null \
	|| docker run \
		--rm \
		--detach=true \
		--name $(NAMESPACE)-postgrest \
		--network=$(NAMESPACE) \
		--env PGRST_DB_URI="postgres://postgres@$(NAMESPACE)-userdb/userdb" \
  		--env PGRST_JWT_SECRET="a_test_only_postgrest_jwt_secret" \
  		--env PGRST_DB_ANON_ROLE="api-anon" \
  		--env PGRST_DB_SCHEMA="v1" \
		local/$(NAMESPACE)-postgrest

.PHONY: docker-start-prod
docker-start-prod:
	docker network inspect $(NAMESPACE) \
	|| docker network create $(NAMESPACE)
	docker inspect -f '{{.State.Running}}' $(NAMESPACE) 2>/dev/null \
	|| docker run \
		--detach=true \
		--name $(NAMESPACE) \
		--network=$(NAMESPACE) \
		--env API_URL="https://userdb.hashbang.sh/v1" \
		--expose="2222" \
		-p "2222:2222" \
		local/$(NAMESPACE)

.PHONY: docker-stop
docker-stop:
	docker inspect -f '{{.State.Running}}' $(NAMESPACE) 2>/dev/null \
	&& docker rm -f $(NAMESPACE) || true
	docker inspect -f '{{.State.Running}}' $(NAMESPACE)-userdb 2>/dev/null \
	&& docker rm -f $(NAMESPACE)-userdb || true
	docker inspect -f '{{.State.Running}}' $(NAMESPACE)-postgrest 2>/dev/null \
	&& docker rm -f $(NAMESPACE)-postgrest || true

.PHONY: docker-logs
docker-logs:
	scripts/docker-logs $(NAMESPACE) $(NAMESPACE)-userdb $(NAMESPACE)-postgrest

.PHONY: docker-clean
docker-clean: docker-stop
	docker image rm local/$(NAMESPACE)
	docker image rm local/$(NAMESPACE)-test
	docker image rm local/$(NAMESPACE)-postgrest
	docker image rm local/$(NAMESPACE)-userdb

.PHONY: docker-test
docker-test:
	docker run \
		-it \
		--rm \
		--hostname=$(NAMESPACE)-test \
		--name $(NAMESPACE)-test \
		--network=$(NAMESPACE) \
		--env CONTAINER=$(NAMESPACE) \
		--env PGPASSWORD=test_password \
		--env PGHOST=$(NAMESPACE)-userdb \
		--env PGDATABASE=userdb \
		--env PGUSER=postgres \
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
		--env CONTAINER=$(NAMESPACE) \
		--env PGPASSWORD=test_password \
		--env PGHOST=$(NAMESPACE)-userdb \
		--env PGDATABASE=userdb \
		--env PGUSER=postgres \
		local/$(NAMESPACE)-test \
		bash

.PHONY: docker-build-test
docker-build-test:
	docker build -t local/$(NAMESPACE)-test test/
