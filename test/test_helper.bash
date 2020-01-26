#!/bin/bash

setup(){
	echo "Settting up test"
}

teardown(){
	echo "Tearing down test"
    psql -c "delete from passwd;";
}

ssh_command(){
	key="${1:-}"
	user="${2:-jdoe}"
	cmd="${3:-}"
	ssh \
		-p 2222 \
		-t \
		-a \
		$([[ "$key" ]] && echo "-i ${HOME}/keys/id_${key}") \
		-o UserKnownHostsFile=/dev/null \
		-o StrictHostKeyChecking=no \
		-o LogLevel=ERROR \
		"$user"@"${CONTAINER}" \
		"$cmd"
}

base64_url_encode(){
	data=${1?}
	echo -n "${data}" \
	| openssl base64 -e -A \
	| sed 's/\+/-/g' \
	| sed 's/\//_/g' \
	| sed -E 's/=+$//'
}

jwt_sig(){
	data=${1?}
	secret=${2?}
	signature=$( \
		echo -n "${data}" \
		| openssl dgst -sha256 -hmac "${secret}" -binary \
		| openssl base64 -e -A \
		| sed 's/\+/-/g' \
		| sed 's/\//_/g' \
		| sed -E 's/=+$//'
	)
	echo -n "${data}"."${signature}"
}

jwt_token(){
	role=${1:-role}
	secret=${2:-a_test_only_postgrest_jwt_secret}
	header="$(base64_url_encode '{"alg":"HS256"}')"
	payload="$(base64_url_encode '{"role":"'"${role}"'"}')"
	echo -n "$(jwt_sig "${header}.${payload}" "${secret}")"
}
