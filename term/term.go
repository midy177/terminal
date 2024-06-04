package term

import (
	"io"
	"os"
	"os/exec"
)

type Terminal struct {
	pty      io.Closer
	cmd      *exec.Cmd
	startDir string
}

func (t *Terminal) startingDir() string {
	if t.startDir == "" {
		home, err := os.UserHomeDir()
		if err == nil {
			return home
		}
	}

	return t.startDir
}
