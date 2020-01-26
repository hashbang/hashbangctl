# hashbangctl
<https://github.com/hashbang/hashbangctl>

## About ##

This dameon allows users to create/manage shell accounts over ssh.

## Features ##

* Current
  * New users get a form to create an account
    * Prefills form with hosts, username and ssh key
      * Randomly populate hosts dropdown from userdb
      * Direct users to run ssh-keygen if no key detected
      * Suggest available username based on incoming username
  * Connection IP rate limiting
  * Basic end to end test suite
* Future
  * If incoming ssh key detected and account exists, direct to management menu
    * Allow users to change their details in UserDB, manage keys, etc
  * Support non-interactive use
  * text captcha
  * k8s deployment boilerplate
    * strict pod security policy
    * strict apparmor/seccomp rules

## Requirements ##
- Docker 19+
- Go 1.13+

## Usage ##

### Build

Build or rebuild the binaries.

```
make
```

### Serve

Start the server

```
make serve
```

### Connect

Connect to the server

```
make connect
```

### Test

Run the test suite

```
make test
```

### Test Shell

Launch a shell in the test suite to run any tests by hand

```
make test-shell
> ssh_command "ed25519" "jdoe" "some-command"
```

### Clean

Stop all containers and cleanup

```
make clean
```
