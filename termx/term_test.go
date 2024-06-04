package termx

import (
	"bufio"
	"io"
	"os"
	"testing"
)

func TestTerm(t *testing.T) {
	cShells := GetShells()
	if len(cShells) == 0 {
		t.Fatal("no shells found")
	}
	cPty, err := NewPTY(&cShells[1])
	if err != nil {
		t.Fatal(err)
	}
	defer cPty.Close()
	// 使用 bufio 包装
	stdinReader := bufio.NewReader(os.Stdin)
	stdoutWriter := bufio.NewWriter(os.Stdout)
	stdErrWriter := bufio.NewWriter(os.Stderr)
	go io.Copy(stdErrWriter, cPty)
	go io.Copy(stdoutWriter, cPty)
	go io.Copy(cPty, stdinReader)
	select {}
}
