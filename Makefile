## Usage

.PHONY: build
build:
	docker build -t local/hashbangctl .

.PHONY: connect
connect: start
	ssh \
		-o UserKnownHostsFile=/dev/null \
		-o StrictHostKeyChecking=no \
		-p2222 localhost

## Development

.PHONY: start
start:
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

.PHONY: stop
stop:
	docker inspect -f '{{.State.Running}}' hashbangctl 2>/dev/null \
	&& docker rm -f hashbangctl || true

.PHONY: log
log:
	docker logs -f hashbangctl

.PHONY: clean
clean: stop
	docker image rm local/hashbangctl

## Testing

.PHONY: test-ssh
test-ssh: start
	SSH_AUTH_SOCK="" \
	ssh \
		-i test/keys/id_ed25519 \
		-o UserKnownHostsFile=/dev/null \
		-o StrictHostKeyChecking=no \
		-p2222 localhost



.PHONY: test
test: stop start build-test
	docker run \
		--rm \
		--hostname=hashbangctl-test \
		--name hashbangctl-test \
		--network=hashbangctl \
		--env="CONTAINER=hashbangctl" \
		local/hashbangctl-test

.PHONY: test-shell
test-shell: stop start build-test
	docker run \
		--rm \
		-it \
		--hostname=hashbangctl-test \
		--name hashbangctl-test \
		--network=hashbangctl \
		--env CONTAINER="hashbangctl" \
		local/hashbangctl-test \
		bash

.PHONY: build-test
build-test:
	docker build -t local/hashbangctl-test test/
