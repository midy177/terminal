//go:build !windows

package termx

import (
	"errors"
	"github.com/creack/pty"
	"github.com/pkg/sftp"
	"io"
	"os"
	"os/exec"
	"sync/atomic"
)

type unixPty struct {
	pty    *os.File
	closed *atomic.Bool
}

func (t *unixPty) Sftp() (*sftp.Client, error) {
	//TODO implement me
	return nil, errors.New("sftp pty not supported")
}

// CloseSftp close sftp client
func (s *unixPty) CloseSftp() error {
	return nil
}

func (t *unixPty) Resize(rows, cols int) error {
	return pty.Setsize(t.pty, &pty.Winsize{
		Rows: uint16(rows), Cols: uint16(cols),
		//X: uint16(t.Size().Width * scale), Y: uint16(t.Size().Height * scale),
	})
}

func (t *unixPty) Read(p []byte) (n int, err error) {
	if t.closed.Load() {
		return 0, io.EOF
	}
	return t.pty.Read(p)
}

func (t *unixPty) Write(p []byte) (n int, err error) {
	if t.closed.Load() {
		return 0, io.EOF
	}
	return t.pty.Write(p)
}

func (t *unixPty) Close() error {
	if t.closed.CompareAndSwap(false, true) {
		return t.pty.Close()
	}
	return nil
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
		pty:    uPty,
		closed: &atomic.Bool{},
	}, err
}
