package main

import (
	"fmt"
	"log"
	"os"
)

func main() {

	fd := os.NewFile(3, "/proc/self/fd/3")
	defer fd.Close()
	logger := log.New(fd, "", log.Ldate|log.Ltime)
    key := os.Getenv("KEY")
	if key == "none" {
		fmt.Fprintln(
			os.Stderr,
			"\nError: Public key authentication required\n",
			"\nFor help generating a key try:\n",
			"\n$ ssh-keygen -t ed25519 -f \"$HOME/.ssh/id_ed25519\"\n",
		)
		os.Exit(1)
	}
    keys, err := getKeys(key)
	if err != nil {
	    fmt.Fprintln(os.Stderr, "\nError: Unable to get keys list")
	    fmt.Fprintln(os.Stderr, err)
    }
    if len(keys) == 0 {
	    hosts, err := getHosts()
	    if err != nil {
	    	fmt.Fprintln(os.Stderr, "\nError: Unable to get host list")
	    	fmt.Fprintln(os.Stderr, err)
	    	os.Exit(1)
	    }
        createForm(logger, hosts)
    } else {
	    fmt.Fprintln(os.Stderr, "\nError: User with this key already exists")
        users, err := getUsersById(keys[0].Uid)
	    if err != nil {
	    	fmt.Fprintln(os.Stderr, "\nError: Unable to get user list")
	    	fmt.Fprintln(os.Stderr, err)
	    	os.Exit(1)
	    }
        fmt.Println(users)
    }
}
