//go:build windows

package termx

import (
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/ActiveState/termtest/conpty"
)

func (t *Terminal) updatePTYSize(rows, cols int) {
	if t.pty == nil { // during load
		return
	}
	_ = t.pty.(*conpty.ConPty).Resize(uint16(cols), uint16(rows))
}

func (t *Terminal) startPTY() (io.WriteCloser, io.Reader, io.Closer, error) {
	cpty, err := conpty.New(80, 25)
	if err != nil {
		return nil, nil, nil, err
	}

	pid, _, err := cpty.Spawn(
		"C:\\WINDOWS\\System32\\WindowsPowerShell\\v1.0\\powershell.exe",
		[]string{},
		&syscall.ProcAttr{
			Env: os.Environ(),
		},
	)
	if err != nil {
		return nil, nil, nil, err
	}

	t.cmd = &exec.Cmd{}
	process, err := os.FindProcess(pid)
	if err != nil {
		return nil, nil, nil, err
	}
	go func() {
		ps, err := process.Wait()
		if err != nil {
			log.Fatalf("Error waiting for process: %v", err)
		}
		t.cmd.ProcessState = ps
		if t.pty != nil {
			t.pty = nil
			_ = cpty.Close()
		}
	}()

	return cpty.InPipe(), cpty.OutPipe(), cpty, nil
}
