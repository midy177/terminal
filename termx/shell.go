package termx

import (
	"os"
	"sync"
)

type SystemShell struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Command string   `json:"command"`
	Args    []string `json:"args"`
	Env     []string `json:"env"`
	Cwd     string   `json:"cwd"`
	Icon    string   `json:"icon"`
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
