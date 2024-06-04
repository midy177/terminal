package termx

import (
	"os"
	"sync"
)

type SystemShell struct {
	ID      string
	Name    string
	Command string
	Args    []string
	Env     []string
	Cwd     string
	Icon    string
}

var shells []SystemShell
var once sync.Once

func GetShells() []SystemShell {
	once.Do(func() {
		getShells()
	})
	return shells
}

func startDir() string {
	home, err := os.UserHomeDir()
	if err == nil {
		return home
	}
	return ""
}
