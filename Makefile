.PHONY: build
build:
	docker build -t local/hashbangctl .

.PHONY: start
start:
	docker inspect -f '{{.State.Running}}' hashbangctl 2>/dev/null || \
	docker run \
		--detach=true \
		--name hashbangctl \
		-p 2222:2222 \
		local/hashbangctl

.PHONY: stop
stop:
	docker rm -f hashbangctl

.PHONY: shell
shell: start
	while sleep 1; do \
		ssh \
			-o UserKnownHostsFile=/dev/null \
			-o StrictHostKeyChecking=no \
			-p2222 localhost; break;\
	done

.PHONY: logs
logs:
	docker logs hashbangctl

.PHONY: test
test: shell
