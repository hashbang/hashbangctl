package main

import (
    "log"
    "net"
    "io"
    "sync"
    "fmt"
    "strings"
    "syscall"
    "unsafe"
    "os"
    "os/exec"
    "path/filepath"
    "crypto/rsa"
    "crypto/rand"
    "encoding/binary"
    "golang.org/x/crypto/ssh"
    "github.com/kr/pty"
)

func setWinsize(fd uintptr, w, h uint32) {
	syscall.Syscall(
        syscall.SYS_IOCTL,
        fd,
        uintptr(syscall.TIOCSWINSZ),
        uintptr(unsafe.Pointer(
            &struct{ h, w, x, y uint16 }{uint16(h), uint16(w), 0, 0},
        )),
    )
}

func parseDims(b []byte) (uint32, uint32) {
	w := binary.BigEndian.Uint32(b)
	h := binary.BigEndian.Uint32(b[4:])
	return w, h
}

func PtyRun(c *exec.Cmd, tty *os.File) (err error) {
    defer tty.Close()
    c.Stdout = tty
    c.Stdin = tty
    c.Stderr = tty
    c.SysProcAttr = &syscall.SysProcAttr{
        Setctty: true,
        Setsid:  true,
    }
    return c.Start()
}

func main() {

    runDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

    // Setup SSH configuration
    config := ssh.ServerConfig{
        NoClientAuth: true,
        PublicKeyCallback: func(
            conn ssh.ConnMetadata,
            key ssh.PublicKey,
        ) (*ssh.Permissions, error) {
            user := conn.User()
            publicKey := strings.TrimSpace(
                string(ssh.MarshalAuthorizedKey(key)),
            )
            log.Printf("New connection: %s %s", publicKey, user)
            return nil, nil
        },
        KeyboardInteractiveCallback: func(
            ssh.ConnMetadata,
            ssh.KeyboardInteractiveChallenge,
        ) (*ssh.Permissions, error) {
            return nil, nil
        },
    }

    // Setup Host key
    key, _ := rsa.GenerateKey(rand.Reader, 2048)
    signer, _ := ssh.NewSignerFromKey(key)
    config.AddHostKey(signer)

    // Listen
    listener, err := net.Listen("tcp4", ":2222")
    if err != nil { panic(err) }
    for {
        conn, _ := listener.Accept()
        sshConn, chans, _, _ := ssh.NewServerConn(conn, &config)
        log.Println("->", sshConn.RemoteAddr())
        go func(chans <-chan ssh.NewChannel){
            for newChannel := range chans {
                if t := newChannel.ChannelType(); t != "session" {
                    newChannel.Reject(
                        ssh.UnknownChannelType,
                        fmt.Sprintf("unknown channel type: %s", t),
                    );
                    continue;
                }
                channel, requests, _ := newChannel.Accept()
                f, tty, _ := pty.Open()
                go func(in <-chan *ssh.Request) {
                    for req := range in {
                        switch req.Type {
                            case "shell":
                                    filepath.Abs(filepath.Dir(os.Args[0]))
                                cmd := exec.Command(
                                    fmt.Sprintf("%s/client", runDir),
                                )
                                cmd.Env = []string{"TERM=xterm"}
                                err := PtyRun(cmd, tty)
                                if err != nil {
                                    log.Printf("%s", err)
                                }
                                var once sync.Once
                                close := func() {
                                    channel.Close()
                                    log.Println("<-", sshConn.RemoteAddr())
                                }
                                go func() {
                                    io.Copy(channel, f)
                                    once.Do(close)
                                }()
                                go func() {
                                    io.Copy(f, channel)
                                    once.Do(close)
                                }()
                            case "pty-req":
                                termLen := req.Payload[3]
                                w, h := parseDims(req.Payload[termLen+4:])
                                setWinsize(f.Fd(), w, h)
                            case "window-change":
                                w, h := parseDims(req.Payload)
                                setWinsize(f.Fd(), w, h)
                                continue
                        }
                        req.Reply(true, nil)
                    }
                }(requests)
            }
        }(chans)
    }
}
