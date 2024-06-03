package term

import (
	"io"
	"os/exec"
)

type Terminal struct {
	pty io.Closer
	cmd *exec.Cmd
}
