//go:build !windows

package termx

import (
	"github.com/creack/pty"
	"io"
	"os"
	"os/exec"
)

func (t *Terminal) updatePTYSize(rows, cols int) {
	if t.pty == nil { // SSH or other direct connection?
		return
	}
	//scale := float32(1.0)

	_ = pty.Setsize(t.pty.(*os.File), &pty.Winsize{
		Rows: uint16(rows), Cols: uint16(cols),
		//X: uint16(t.Size().Width * scale), Y: uint16(t.Size().Height * scale),
	})
}

func (t *Terminal) startPTY() (io.WriteCloser, io.Reader, io.Closer, error) {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "bash"
	}

	_ = os.Chdir(t.startingDir())
	env := os.Environ()
	env = append(env, "TERM=xterm-256color")
	c := exec.Command(shell)
	c.Env = env
	t.cmd = c

	// Start the command with a pty.
	f, err := pty.Start(c)
	return f, f, f, err
}
