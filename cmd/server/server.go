package main

import (
    "log"
    "io"
    "sync"
    "fmt"
    "net"
    "strings"
    "syscall"
    "unsafe"
    "errors"
    "os"
    "os/exec"
    "path/filepath"
    "encoding/json"
    "encoding/binary"
    "golang.org/x/crypto/ssh"
    "github.com/kr/pty"
)

type sessionInfo struct {
    User          string   `json:"user"`
    ClientVersion string   `json:"clientVersion"`
    RemoteAddr    string   `json:"remoteAddr"`
    Keys          []string `json:"keys"`
}

type Server struct {
    sshConfig   *ssh.ServerConfig
    sessionInfo map[string]sessionInfo
    mu          sync.RWMutex
}

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

func (s *Server) Handle(nConn net.Conn) {
    conn, chans, _, _ := ssh.NewServerConn(nConn, s.sshConfig)
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
                s.mu.RLock()
                si := s.sessionInfo[string(conn.SessionID())]
                s.mu.RUnlock()
                sessionJSON, _ := json.Marshal(si)
                for req := range in {
                    switch req.Type {
                        case "shell":
                            log.Println("->", string(sessionJSON))
                            runDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
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
                                log.Println("<-", string(sessionJSON))
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

func (s *Server) PublicKeyCallback(
    conn ssh.ConnMetadata,
    key ssh.PublicKey,
) (*ssh.Permissions, error) {
    s.mu.Lock()
    si := s.sessionInfo[string(conn.SessionID())]
    si.User = conn.User()
    si.RemoteAddr = fmt.Sprintf("%s", conn.RemoteAddr())
    si.ClientVersion = fmt.Sprintf("%s", conn.ClientVersion())
    si.Keys = append(
        si.Keys,
        strings.TrimSpace(
            string(ssh.MarshalAuthorizedKey(key)),
        ),
    )
    s.sessionInfo[string(conn.SessionID())] = si
    s.mu.Unlock()

    return nil, errors.New("")
}

func (s *Server) KeyboardInteractiveCallback(
    ssh.ConnMetadata,
    ssh.KeyboardInteractiveChallenge,
) (*ssh.Permissions, error) {
    return nil, nil
}
