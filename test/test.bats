load test_helper

@test "Cannot login without an ssh key" {
	run ssh_command
	[ "$status" -eq 255 ]
}

@test "Can login with an ed25519 ssh key" {
	run ssh_command "ed25519" "jdoe" "debug"
	[ "$status" -eq 0 ]
	[[ "$output" == *"ssh-ed25519"* ]]
}

@test "Cannot run an invalid command" {
	run ssh_command "ed25519" "jdoe" "invalid"
	[ "$status" -eq 255 ]
	[[ "$output" == *"Unknown command"* ]]
}

@test "Can create user with an ed25519 ssh key" {
	run tmux_command "ssh_command ed25519"
	run tmux_keys TAB TAB TAB ENTER
	[ "$status" -eq 0 ]
}
