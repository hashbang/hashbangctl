# hashbangctl
<https://github.com/hashbang/hashbangctl>

## About ##

This dameon allows users to create/manage shell accounts over ssh.

## Features ##

* Current
  * Users can ssh in and get a form to create a user (does nothing yet)
  * If incoming ssh key not detected, direct users to run ssh-keygen
  * Prefills form with username and ssh key from incoming connection
* Future
  * Will randomly set a server but allow users to change it
  * If incoming ssh key detected and account exists, direct to management menu
    * Allow users to change their details in UserDB, manage keys, etc
  * Abuse mitigation
    * text captcha
    * rate limiting
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
