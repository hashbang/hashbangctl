load test_helper

@test "Can connect to userdb PostgreSQL" {
    sleep 5
	run pg_isready -U postgres -h hashbangctl-userdb;
	[ "$status" -eq 0 ]
	echo "$output" | grep "accepting connections"
}

@test "Can connect to userdb PostgREST" {
	run curl http://hashbangctl-postgrest:3000
	[ "$status" -eq 0 ]
	echo "$output" | grep "swagger"
}

@test "Cannot create user anonymously via PostgREST" {
	run curl http://hashbangctl-postgrest:3000/passwd \
		-H "Content-Type:application/json" \
		-X POST \
		--data-binary @- <<-EOF
			{
				"name": "testuser",
				"host": "test.hashbang.sh",
				"data": {
					"shell": "/bin/bash",
					"ssh_keys": ["$(cat keys/id_ed25519.pub)"]
				}
			}
			EOF
	[ "$status" -eq 0 ]
	echo "$output" | grep "permission denied"
}

@test "Can create user with a valid host and valid auth via PostgREST" {

	run curl http://hashbangctl-postgrest:3000/passwd \
		-H "Content-Type: application/json" \
		-H "Authorization: Bearer $(jwt_token 'api-user-create')" \
		-X POST \
		--data-binary @- <<-EOF
			{
				"name": "testuser42",
				"host": "te1.hashbang.sh",
				"data": {
					"shell": "/bin/bash",
					"ssh_keys": ["$(cat keys/id_ed25519.pub)"]
				}
			}
			EOF
	[ "$status" -eq 0 ]

	run curl http://hashbangctl-postgrest:3000/passwd?name=eq.testuser42
	echo "$output" | grep "testuser42"
}

@test "Cannot login without an ssh key" {
	run ssh_command
	[ "$status" -eq 255 ]
}

#@test "Can login with an ed25519 ssh key" {
#	run ssh_command "ed25519"
#	[ "$status" -eq 0 ]
#}
