//go:build !windows

package termx

import (
	"context"
	"errors"
	"fmt"
	"github.com/creack/pty"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"os/exec"
	"sync/atomic"
)

type unixPty struct {
	cmd    *exec.Cmd
	pty    *os.File
	closed *atomic.Bool
	cancel func()
}

func (t *unixPty) Ssh() (*ssh.Client, error) {
	//TODO implement me
	return nil, errors.New("creack pty not supported")
}

func (t *unixPty) Sftp() (*sftp.Client, error) {
	//TODO implement me
	return nil, errors.New("sftp pty not supported")
}

// CloseSftp close sftp client
func (t *unixPty) CloseSftp() error {
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
		t.cancel()
		if sf, ok := t.cmd.Stdout.(*os.File); ok {
			_ = sf.Close()
		}
		return t.pty.Close()
	}
	return nil
}

func NewPTY(s *SystemShell) (PtyX, error) {
	env := os.Environ()
	env = append(env, s.Env...)
	//c := exec.Command(s.Command, s.Args...)
	ctx, cancel := context.WithCancel(context.Background())
	c := exec.CommandContext(ctx, s.Command, s.Args...)
	c.Env = env

	if s.Cwd != "" {
		c.Dir = s.Cwd
	}
	// Start the command with a pty.
	uPty, err := pty.Start(c)
	closed := &atomic.Bool{}
	go func() {
		err = c.Wait()
		if err != nil {
			fmt.Printf("pty shell exited: %s\n", err)
		}
		if closed.CompareAndSwap(false, true) {
			cancel()
			_ = uPty.Close()
			if sf, ok := c.Stdout.(*os.File); ok {
				_ = sf.Close()
			}
		}
	}()

	return &unixPty{
		cmd:    c,
		pty:    uPty,
		closed: closed,
		cancel: cancel,
	}, err
}
