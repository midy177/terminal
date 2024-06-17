package termx

import (
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/term"
	"io"
	"log"
	"os"
	"testing"
)

func TestTerm(t *testing.T) {
	termFD := int(os.Stdin.Fd())
	termState, _ := terminal.MakeRaw(termFD)
	defer terminal.Restore(termFD, termState)
	termWidth, termHeight, _ := term.GetSize(termFD)
	cPty, err := NewSshPTY("milesight",
		"CpE*!Vy4g6B@",
		"192.168.5.135",
		22, nil,
		termHeight, termWidth)
	if err != nil {
		log.Fatal(err)
	}
	defer cPty.Close()
	// 使用 bufio 包装
	//stdinReader := bufio.NewReader(os.Stdin)
	//stdoutWriter := bufio.NewWriter(os.Stdout)
	//stdErrWriter := bufio.NewWriter(os.Stderr)
	go io.Copy(os.Stdout, cPty)
	go io.Copy(os.Stderr, cPty)
	go io.Copy(cPty, os.Stdin)
	select {}
}
