#!/bin/bash

setup(){
    echo "setup"
}

teardown(){
    echo "teardown"
}

ssh_command(){
	key="${1:-}"
	user="${2:-jdoe}"
	cmd="${3:-}"
	ssh \
		-p 2222 \
		-a \
		$([[ "$key" ]] && echo "-i ${HOME}/keys/id_${key}") \
		-o UserKnownHostsFile=/dev/null \
		-o StrictHostKeyChecking=no \
		-o LogLevel=ERROR \
		"$user"@"${CONTAINER}" \
		"$cmd"
}
