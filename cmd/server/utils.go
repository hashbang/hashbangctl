package main

import (
	"encoding/binary"
	"golang.org/x/crypto/ssh"
	"os"
	"os/exec"
	"strings"
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
