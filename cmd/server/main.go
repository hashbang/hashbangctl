package main

import (
    "log"
    "net"
    "crypto/rsa"
    "crypto/rand"
    "golang.org/x/crypto/ssh"
)

func main() {

    key, _ := rsa.GenerateKey(rand.Reader, 2048)
    signer, _ := ssh.NewSignerFromKey(key)

    server := &Server{
		sessionInfo:  make(map[string]sessionInfo),
	}
	server.sshConfig = &ssh.ServerConfig{
		KeyboardInteractiveCallback: server.KeyboardInteractiveCallback,
		PublicKeyCallback: server.PublicKeyCallback,
	}
    server.sshConfig.AddHostKey(signer)

    listener, err := net.Listen("tcp4", ":2222")
    if err != nil { panic(err) }
    for {
        conn, err := listener.Accept()
        if err != nil { log.Println("!! %s", err); continue; }
        go server.Handle(conn)
    }

}
