package termx

import (
	"io"
	"os"
	"testing"
)

func TestTerm(t *testing.T) {
	cShells := GetShells()
	if len(cShells) == 0 {
		t.Fatal("no shells found")
	}
	cPty, err := NewPTY(&cShells[0])
	if err != nil {
		t.Fatal(err)
	}
	defer cPty.Close()

	go io.Copy(os.Stdout, cPty)
	go io.Copy(cPty, os.Stdin)
	select {}
}
