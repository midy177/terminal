//go:build !windows

package termx

import (
	"github.com/creack/pty"
	"os"
	"os/exec"
)

type unixPty struct {
	pty *os.File
}

func (t *unixPty) Resize(rows, cols int) error {
	return pty.Setsize(t.pty, &pty.Winsize{
		Rows: uint16(rows), Cols: uint16(cols),
		//X: uint16(t.Size().Width * scale), Y: uint16(t.Size().Height * scale),
	})
}

func (t *unixPty) Read(p []byte) (n int, err error) {
	return t.pty.Read(p)
}

func (t *unixPty) Write(p []byte) (n int, err error) {
	return t.pty.Write(p)
}

func (t *unixPty) Close() error {
	return t.pty.Close()
}

func NewPTY(s *SystemShell) (PtyX, error) {
	env := os.Environ()
	env = append(env, s.Env...)
	c := exec.Command(s.Command, s.Args...)
	c.Env = env

	if s.Cwd != "" {
		c.Dir = s.Cwd
	}
	// Start the command with a pty.
	uPty, err := pty.Start(c)
	return &unixPty{
		pty: uPty,
	}, err
}
