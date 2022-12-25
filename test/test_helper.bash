#!/bin/bash

setup(){
	echo "Settting up test"
	psql -c "insert into hosts (name,maxusers) values ('test.hashbang.sh','500');";
	psql -c "insert into hosts (name,maxusers) values ('test2.hashbang.sh','500');";
	tmux new -d -s hashbangctl-test
	tmux send-keys -t hashbangctl-test "source test/test_helper.bash" ENTER
}

teardown(){
	echo "Tearing down test"
	psql -c "delete from passwd;";
	psql -c "delete from hosts;";
	tmux kill-session -t hashbangctl-test
}

tmux_commmand(){
	cmd="${1:-}"
	tmux send-keys -t hashbangctl-test "$cmd"
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
