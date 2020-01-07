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
  * k8s deployment boilerplate
    * strict pod security policy
    * strict apparmor/seccomp rules

## Requirements ##
- Docker 19+

## Build

```
make
```

## Start

```
make start
```

## Stop

```
make stop
```

## Logs

```
make logs
```

## Test

```
make test
```
