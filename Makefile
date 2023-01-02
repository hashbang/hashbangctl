NAMESPACE ?= hashbangctl

.PHONY: build
build:
	GOBIN=$(PWD)/bin \
	GOPATH=$(PWD)/go \
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64 \
	go install ./...

.PHONY: serve
serve: build
	API_URL="https://userdb.hashbang.sh/v1" \
	API_TOKEN="eyJhbGciOiJIUzI1NiJ9.eyJyb2xlIjoiYXBpLXVzZXItbWFuYWdlIn0.gEKj0rWfh5_cnhfAk3U-oGL1_WN6ZrbywfIwxXx_wDA" \
	HOST_KEY_SEED="This is an insecure seed" \
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
test: docker-build docker-build-test docker-stop docker-start
	docker run \
		-it \
		--rm \
		--hostname=$(NAMESPACE)-test \
		--name $(NAMESPACE)-test \
		--network=userdb \
		--env UID=$(shell id -u) \
		--env GID=$(shell id -g) \
		--env CONTAINER=$(NAMESPACE) \
		--env PGPASSWORD=test_password \
		--env PGHOST=userdb-postgres \
		--env PGDATABASE=postgres \
		--env PGUSER=postgres \
		--volume $(PWD)/test:/home/test \
		local/$(NAMESPACE)-test \
		scripts/docker-shell bats test.bats

.PHONY: test-shell
test-shell: docker-build docker-build-test docker-stop docker-start
	docker run \
		-it \
		--rm \
		--hostname=$(NAMESPACE)-test \
		--name $(NAMESPACE)-test \
		--network=userdb \
		--env UID=$(shell id -u) \
		--env GID=$(shell id -g) \
		--env CONTAINER=$(NAMESPACE) \
		--env PGPASSWORD=test_password \
		--env PGHOST=userdb-postgres \
		--env PGDATABASE=postgres \
		--env PGUSER=postgres \
		--volume $(PWD)/test:/home/test \
		local/$(NAMESPACE)-test \
		scripts/docker-shell bash

.PHONY: clean
clean: docker-clean
	rm -rf ./go ./bin

.PHONY: initdb
initdb:
	docker exec --user postgres -it "userdb-postgres" \
		psql -c "delete from passwd; delete from hosts;";
	docker exec --user postgres -it "userdb-postgres" \
		psql -c "insert into hosts (name,maxusers) values ('local1.hashbang.sh','500');";
	docker exec --user postgres -it "userdb-postgres" \
		psql -c "insert into hosts (name,maxusers) values ('local2.hashbang.sh','500');";


.PHONY: docker-logs
docker-logs:
	scripts/docker-logs $(NAMESPACE) userdb-postgres userdb-postgrest

.PHONY: docker-build
docker-build:
	docker build -t local/$(NAMESPACE) .

.PHONY: docker-build-test
docker-build-test:
	docker build -t local/$(NAMESPACE)-test test/

.PHONY: docker-start
docker-start: docker-build
	$(MAKE) -C test/modules/userdb docker-start
	docker inspect -f '{{.State.Running}}' $(NAMESPACE) 2>/dev/null \
	|| docker run \
		--detach=true \
		--name $(NAMESPACE) \
		--network=userdb \
		--env API_URL="http://userdb-postgrest:3000" \
		--env API_TOKEN="eyJhbGciOiJIUzI1NiJ9.eyJyb2xlIjoiYXBpLXVzZXItbWFuYWdlIn0.gEKj0rWfh5_cnhfAk3U-oGL1_WN6ZrbywfIwxXx_wDA" \
		--env HOST_KEY_SEED="This is an insecure seed" \
		--expose="2222" \
		-p "2222:2222" \
		local/$(NAMESPACE)

.PHONY: docker-stop
docker-stop:
	docker inspect -f '{{.State.Running}}' $(NAMESPACE) 2>/dev/null \
	&& docker rm -f $(NAMESPACE) || true
	$(MAKE) -C test/modules/userdb docker-stop
