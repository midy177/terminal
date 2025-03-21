//go:build !windows

package termx

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/creack/pty"
)

func supportLogin() bool {
	user := os.Getenv("USER")
	if len(user) == 0 {
		return false
	}
	bash := exec.Command(
		"sudo",
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
	_, err = tty.Read(buf)
	if err != nil {
		fmt.Println("not support use 'login -f " + user + "' to start pty,err: " + err.Error())
		return false
	}
	//fmt.Println("support use 'login -f " + user + "' to start pty,stdout: " + string(buf[:read]))
	return true
}

func getShortHostname() string {
	cmd := exec.Command("hostname", "-s")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error run(hostname -s):", err)
		return ""
	}
	// 将输出转换为字符串并去除可能存在的空白字符
	hostname := strings.TrimSpace(string(output))
	//fmt.Println("Short hostname:", hostname)
	return hostname
}

func getShells() {
	sbl := supportLogin()
	term := os.Getenv("xterm")
	if term == "" {
		term = "xterm-256color"
	}
	home := startDir()
	user := os.Getenv("USER")
	hostname := getShortHostname()
	title := user
	if hostname != "" {
		title += "@" + hostname
	}
	if sbl && len(user) > 0 {
		shells = append(shells, SystemShell{
			//ID:      "login",
			Name:    title,
			Command: "sudo",
			Args:    []string{"login", "-f", user},
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
			//fmt.Println(line)
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
