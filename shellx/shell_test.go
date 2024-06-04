package shellx

import (
	"fmt"
	"github.com/containerd/console"
	"testing"
)

func TestName(t *testing.T) {
	current := console.Current()
	defer current.Reset()

	if err := current.SetRaw(); err != nil {
	}
	ws, err := current.Size()
	if err != nil {
		fmt.Println(err)
	}
	current.Resize(ws)
}
