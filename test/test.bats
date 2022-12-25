load test_helper

@test "Cannot login without an ssh key" {
	run ssh_command
	[ "$status" -eq 255 ]
}

@test "Can login with an ed25519 ssh key" {
	run ssh_command "ed25519" "jdoe" "version"
	[ "$status" -eq 0 ]
}

@test "Can login with an ed25519 ssh key" {
	run ssh_command "ed25519"
	[ "$status" -eq 0 ]
}

@test "Can create user with an ed25519 ssh key" {
	run tmux_command 'ssh_command ed25519'
	run tmux_command 'TAB TAB TAB ENTER'
	[ "$status" -eq 0 ]
}
