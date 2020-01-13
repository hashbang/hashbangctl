load test_helper

@test "Can login with an ed25519 ssh key" {
    run ssh_command "ed25519"
    [ "$status" -eq 0 ]
}

@test "Can't login without an ssh key" {
    run ssh_command ""
    [ "$status" -eq 1 ]
}
