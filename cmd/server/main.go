package main

import (
    "io"
    "fmt"
    "os"
    "os/exec"
    "syscall"
    "unsafe"
    "log"
    "github.com/gliderlabs/ssh"
    "github.com/kr/pty"
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

    log.Println("starting ssh server on port 2222...")
    log.Fatal(ssh.ListenAndServe(":2222", nil))
}
