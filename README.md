# hashbangctl
<https://github.com/hashbang/hashbangctl>

## About ##

This dameon allows users to create/manage shell accounts over ssh.

## Features ##

* Current
  * Users can ssh in and get a form to create a user (does nothing yet)
* Future
  * Will prefill form with username and ssh keys from incoming connection
  * Will randomly set a server but allow users to change it
  * If incoming ssh key not detected, direct users to run our key setup script
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

## Usage ##

### Build

Build or Rebuild the hashbangctl container

```
make build
```

### Connect

Connect to the hashbangctl container using your own ssh credentials

```
make connect
```

## Develop

### Start

Start the hashbangctl docker container

```
make start
```

### Stop

Stop the hashbangctl docker container

```
make stop
```

### Log

Tail logs from the hashbangctl docker container

```
make log
```

### Clean

Stop all containers and cleanup

```
make clean
```

## Testing

### Test

Run the BATS test suite

```
make test
```

### test-ssh

ssh to the hashbangctl container using test credentials

```
make test-ssh
```

### Test Shell

Launch a shell in the test suite to run any tests by hand

```
make test-shell
> ssh_command "ed25519" "jdoe" "some-command"
```
