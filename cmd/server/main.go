package main

import (
	"encoding/json"
	"crypto/sha512"
	"crypto/ed25519"
	"encoding/binary"
	"math/rand"
	"fmt"
	"github.com/creack/pty"
	"golang.org/x/crypto/ssh"
	"golang.org/x/time/rate"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

var limiter = NewIPRateLimiter(rate.Every(time.Minute), 10)

var (
	hostPrivateKeySigner ssh.Signer
)

type LoginData struct {
	UserName   string `json:"name"`
	IpAddress  string `json:"ip"`
	SshKey     string `json:"key"`
	SshVersion string `json:"version"`
}

func handleConnection(nConn net.Conn, sshConfig *ssh.ServerConfig) {
	conn, chans, _, _ := ssh.NewServerConn(nConn, sshConfig)
	addr, _ := conn.Conn.RemoteAddr().(*net.TCPAddr)
	ipAddress := addr.IP.String()
	sshKey := "none"
	if conn.Permissions != nil {
		if conn.Permissions.Extensions != nil {
			if k, ok := conn.Permissions.Extensions["pubkey"]; ok {
				sshKey = k
			}
		}
	}

	limiter := limiter.GetLimiter(ipAddress)

	loginData := LoginData{
		UserName:   conn.Conn.User(),
		IpAddress:  ipAddress,
		SshKey:     sshKey,
		SshVersion: fmt.Sprintf("%s", conn.Conn.ClientVersion()),
	}
	jsonLoginData, _ := json.Marshal(loginData)

	go func(chans <-chan ssh.NewChannel) {
		for newChannel := range chans {

			if !limiter.Allow() {
				log.Println(
					"[server] !!",
					"Rate Limit exceeded",
					string(jsonLoginData),
				)
				newChannel.Reject(
					ssh.ConnectionFailed,
					fmt.Sprintf("Rate limit exceeded"),
				)
				continue
			}

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

				log.Println("[server] ->", string(jsonLoginData))

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
							log.Println("[server] <-", string(jsonLoginData))
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

type SyncReader struct {
	lk     sync.Mutex
	reader io.Reader
}

func NewSyncReader(reader io.Reader) *SyncReader {
	return &SyncReader{
		reader: reader,
	}
}

func (r *SyncReader) Read(p []byte) (n int, err error) {
	r.lk.Lock()
	defer r.lk.Unlock()
	return r.reader.Read(p)
}

func init() {
	hash := sha512.New()
	io.WriteString(hash, os.Getenv("HOST_KEY_SEED"))
	var seed uint64 = binary.BigEndian.Uint64(hash.Sum(nil)[:8])
	var reader = NewSyncReader(rand.New(rand.NewSource(int64(seed))))
	pubKey, privKey, _ := ed25519.GenerateKey(reader)
	sshPubKey, _ := ssh.NewPublicKey(pubKey)
	hostPrivateKeySigner, _ = ssh.NewSignerFromKey(privKey)
	fmt.Println("SSH Daemon Started")
	fmt.Printf("Host Key: %s", ssh.MarshalAuthorizedKey(sshPubKey))
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
