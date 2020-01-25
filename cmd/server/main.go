package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/kr/pty"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"unsafe"
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
	go func(chans <-chan ssh.NewChannel) {
		for newChannel := range chans {
			if t := newChannel.ChannelType(); t != "session" {
				newChannel.Reject(
					ssh.UnknownChannelType,
					fmt.Sprintf("unknown channel type: %s", t),
				)
				continue
			}
			channel, requests, _ := newChannel.Accept()
			f, tty, _ := pty.Open()
			go func(in <-chan *ssh.Request) {
				sshKey := "none"
				if conn.Permissions != nil {
					if conn.Permissions.Extensions != nil {
						if k, ok := conn.Permissions.Extensions["pubkey"]; ok {
							sshKey = k
						}
					}
				}

				loginData := LoginData{
					UserName:   conn.Conn.User(),
					IpAddress:  fmt.Sprintf("%s", conn.Conn.RemoteAddr()),
					SshVersion: fmt.Sprintf("%s", conn.Conn.ClientVersion()),
					SshKey:     sshKey,
				}
				jsonLoginData, _ := json.Marshal(loginData)

				log.Println("[server] ++", string(jsonLoginData))
				for req := range in {
					switch req.Type {
					case "shell":
						runDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
						pr, pw, _ := os.Pipe()
						defer pw.Close()
						cmd := exec.Command(
							fmt.Sprintf("%s/client", runDir),
						)
						cmd.Env = append(
							os.Environ(),
							"TERM=xterm-256color",
							fmt.Sprintf("IP=%s", loginData.IpAddress),
							fmt.Sprintf("VERSION=%s", loginData.SshVersion),
							fmt.Sprintf("USER=%s", loginData.UserName),
							fmt.Sprintf("KEY=%s", loginData.SshKey),
						)
						cmd.ExtraFiles = []*os.File{pw}
						err := ptyRun(cmd, tty)
						if err != nil {
							log.Printf("%s", err)
						}
						var once sync.Once

						close := func() {
							channel.Close()
							log.Println("[server] --", string(jsonLoginData))
						}
						go func() {
							io.Copy(channel, f)
							once.Do(close)
						}()
						go func() {
							io.Copy(f, channel)
							once.Do(close)
						}()
						go func() {
							defer pr.Close()
							io.Copy(os.Stdout, pr)
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

type LoginData struct {
	UserName   string `json:"name"`
	IpAddress  string `json:"ip"`
	SshKey     string `json:"key"`
	SshVersion string `json:"version"`
}

func init() {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	hostPrivateKeySigner, _ = ssh.NewSignerFromKey(key)

}

func main() {

	sshConfig := &ssh.ServerConfig{
		KeyboardInteractiveCallback: keyboardInteractiveCallback,
		PublicKeyCallback:           publicKeyCallback,
	}
	sshConfig.AddHostKey(hostPrivateKeySigner)

	listener, err := net.Listen("tcp4", ":2222")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("!! %s", err)
			continue
		}
		go handleConnection(conn, sshConfig)
	}

}
