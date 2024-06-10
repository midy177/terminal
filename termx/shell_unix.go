//go:build !windows

package termx

import (
	"bufio"
	"fmt"
	"github.com/creack/pty"
	"os"
	"os/exec"
	"strings"
)

func supportLogin() bool {
	user := os.Getenv("USER")
	if len(user) == 0 {
		return false
	}
	bash := exec.Command(
		"login",
		"-f",
		user,
	)
	tty, err := pty.Start(bash)
	if err != nil {
		fmt.Println("failed to start pty,err: " + err.Error())
		return false
	}
	defer func() {
		_ = bash.Process.Kill()
		_, _ = bash.Process.Wait()
		_ = bash.Process.Release()
		_ = tty.Close()
	}()
	var buf = make([]byte, 1024)
	read, err := tty.Read(buf)
	if err != nil {
		fmt.Println("not support use 'login -f " + user + "' to start pty,err: " + err.Error())
		return false
	}
	fmt.Println("support use 'login -f " + user + "' to start pty,stdout: " + string(buf[:read]))
	return true
}

func getShells() {
	sbl := supportLogin()
	term := os.Getenv("xterm")
	if term == "" {
		term = "xterm-256color"
	}
	home := startDir()
	user := os.Getenv("USER")
	if sbl && len(user) > 0 {
		shells = append(shells, SystemShell{
			//ID:      "login",
			Name:    "Unix(login)",
			Command: "login",
			Args:    []string{"-f", user},
			Env:     []string{"TERM=" + term},
			Cwd:     home,
			Icon:    "/assets/icons/linux.svg",
		})
	} else {
		file, err := os.Open("/etc/shells")
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "#") {
				continue
			}
			fmt.Println(line)
			parts := strings.Split(line, "/")
			if len(parts) > 0 {
				lastPart := parts[len(parts)-1]
				shells = append(shells, SystemShell{
					//ID:      lastPart,
					Name:    "Unix(" + lastPart + ")",
					Command: line,
					Args:    nil,
					Env:     []string{"TERM=" + term},
					Cwd:     home,
					Icon:    "/assets/icons/linux.svg",
				})
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading file:", err)
		}
	}
}
