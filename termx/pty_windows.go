//go:build windows

package termx

import (
	"github.com/UserExistsError/conpty"
	"strings"
)

type windowsPty struct {
	pty *conpty.ConPty
}

func (t *windowsPty) Resize(rows, cols int) error {
	return t.pty.Resize(rows, cols)
}

func (t *windowsPty) Read(p []byte) (n int, err error) {
	return t.pty.Read(p)
}

func (t *windowsPty) Write(p []byte) (n int, err error) {
	return t.pty.Write(p)
}

func (t *windowsPty) Close() error {
	return t.pty.Close()
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
	return &windowsPty{pty: wPty}, nil
}
