package main

import (
    "io"
    "fmt"
    "os"
    "strings"
    "os/exec"
    "syscall"
    "unsafe"
    "log"
    "github.com/gliderlabs/ssh"
    "github.com/kr/pty"
    gossh "golang.org/x/crypto/ssh"
)

func setWinsize(f *os.File, w, h int) {
	syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), uintptr(syscall.TIOCSWINSZ),
		uintptr(unsafe.Pointer(&struct{ h, w, x, y uint16 }{uint16(h), uint16(w), 0, 0})))
}

func main() {
    ssh.Handle(func(s ssh.Session) {
	    ptyReq, winCh, isPty := s.Pty()
	    if !isPty {
	        io.WriteString(s, "only interactive terminals are supported")
	        s.Exit(1)
	        return
	    }
        cmd := exec.Command("/client")
        cmd.Env = append(cmd.Env, fmt.Sprintf("TERM=%s", ptyReq.Term))

        user := s.User()
        cmd.Env = append(cmd.Env, fmt.Sprintf("USERNAME=%s", ptyReq.Term))

        publicKey := strings.TrimSpace(string(gossh.MarshalAuthorizedKey(s.PublicKey())))
        cmd.Env = append(cmd.Env, fmt.Sprintf("PUBLIC_KEY=%s", publicKey))

        log.Println(fmt.Sprintf("Connection: %s %s", publicKey, user))

	    f, err := pty.Start(cmd)
	    if err != nil {
	    	panic(err)
	    }
	    go func() {
	    	for win := range winCh {
	    		setWinsize(f, win.Width, win.Height)
	    	}
	    }()
	    go func() {
	    	io.Copy(f, s) // stdin
	    }()
	    io.Copy(s, f) // stdout
	    cmd.Wait()

    })
    publicKeyOption := ssh.PublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
		return true
	})

    log.Println("starting ssh server on port 2222...")
    log.Fatal(ssh.ListenAndServe(":2222", nil, publicKeyOption))
}
