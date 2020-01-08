load test_helper

@test "Can not login without an ssh key" {
    run ssh_command ""
    [ "$status" -eq 255 ]
}
