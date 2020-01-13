package main

import (
    "log"
    "net"
    "crypto/rsa"
    "crypto/rand"
    "io"
    "sync"
    "fmt"
    "strings"
    "syscall"
    "unsafe"
    "os"
    "os/exec"
    "path/filepath"
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

func ptyRun(c *exec.Cmd, tty *os.File) (err error) {
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

func publicKeyCallback(
    conn ssh.ConnMetadata,
    key ssh.PublicKey,
) (*ssh.Permissions, error) {
    return &ssh.Permissions{
        Extensions: map[string]string{
            "pubkey": strings.TrimSpace(string(ssh.MarshalAuthorizedKey(key))),
        },
    }, nil
}

func keyboardInteractiveCallback(
    ssh.ConnMetadata,
    ssh.KeyboardInteractiveChallenge,
) (*ssh.Permissions, error) {
    return nil, nil
}

func handleConnection(nConn net.Conn, sshConfig *ssh.ServerConfig) {
    conn, chans, _, _ := ssh.NewServerConn(nConn, sshConfig)
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
                user := conn.Conn.User()
                ip := fmt.Sprintf("%s", conn.Conn.RemoteAddr())
                version := fmt.Sprintf("%s", conn.Conn.ClientVersion())
                key := "none"
                if conn.Permissions != nil {
                    if conn.Permissions.Extensions != nil {
                        if k, ok := conn.Permissions.Extensions["pubkey"]; ok {
                            key = k
                        }
                    }
                }
                log.Println("->",ip,version,user,key)
                for req := range in {
                    switch req.Type {
                        case "shell":
                            runDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
                            cmd := exec.Command(
                                fmt.Sprintf("%s/client", runDir),
                            )
                            cmd.Env = []string{"TERM=xterm"}
                            err := ptyRun(cmd, tty)
                            if err != nil {
                                log.Printf("%s", err)
                            }
                            var once sync.Once
                            close := func() {
                                channel.Close()
                                log.Println("<-",ip,version,user,key)
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

var (
	hostPrivateKeySigner ssh.Signer
)

func init() {
    key, _ := rsa.GenerateKey(rand.Reader, 2048)
    hostPrivateKeySigner, _ = ssh.NewSignerFromKey(key)
}

func main() {

    sshConfig := &ssh.ServerConfig{
		KeyboardInteractiveCallback: keyboardInteractiveCallback,
		PublicKeyCallback: publicKeyCallback,
	}
    sshConfig.AddHostKey(hostPrivateKeySigner)

    listener, err := net.Listen("tcp4", ":2222")
    if err != nil { panic(err) }
    for {
        conn, err := listener.Accept()
        if err != nil { log.Println("!! %s", err); continue; }
        go handleConnection(conn, sshConfig)
    }

}
