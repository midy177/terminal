//go:build windows

package termx

import (
	"context"
	"errors"
	"fmt"
	"github.com/UserExistsError/conpty"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"strings"
	"sync/atomic"
)

type windowPty struct {
	pty    *conpty.ConPty
	closed *atomic.Bool
}

func (t *windowPty) Ssh() (*ssh.Client, error) {
	//TODO implement me
	return nil, errors.New("window conpty not supported")
}

func (t *windowPty) Sftp() (*sftp.Client, error) {
	//TODO implement me
	return nil, errors.New("sftp pty not supported")
}

// CloseSftp close sftp client
func (t *windowPty) CloseSftp() error {
	return nil
}

func (t *windowPty) Resize(rows, cols int) error {
	return t.pty.Resize(cols, rows)
}

func (t *windowPty) Read(p []byte) (n int, err error) {
	if t.closed.Load() {
		return 0, io.EOF
	}
	return t.pty.Read(p)
}

func (t *windowPty) Write(p []byte) (n int, err error) {
	if t.closed.Load() {
		return 0, io.EOF
	}
	return t.pty.Write(p)
}

func (t *windowPty) Close() error {
	if t.closed.CompareAndSwap(false, true) {
		return t.pty.Close()
	}
	return nil
}

func NewPTY(s *SystemShell) (PtyX, error) {
	if !conpty.IsConPtyAvailable() {
		return nil, conpty.ErrConPtyUnsupported
	}
	cmdLine := s.Command + " " + strings.Join(s.Args, " ")
	var options []conpty.ConPtyOption
	options = append(options, conpty.ConPtyDimensions(1000, 700))
	if s.Cwd != "" {
		options = append(options, conpty.ConPtyWorkDir(s.Cwd))
	}
	wPty, err := conpty.Start(cmdLine, options...)
	if err != nil {
		return nil, err
	}
	closed := &atomic.Bool{}
	go func() {
		exitCode, err := wPty.Wait(context.Background())
		if err != nil {
			fmt.Printf("conpty is exiting with err: %s\n", err)
		} else {
			fmt.Println("conpty is exiting with exitCode: ", exitCode)
		}
		if closed.CompareAndSwap(false, true) {
			_ = wPty.Close()
		}
	}()

	return &windowPty{
		pty:    wPty,
		closed: closed,
	}, nil
}
