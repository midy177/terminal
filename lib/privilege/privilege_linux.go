//go:build linux
// +build linux

package privilege

import (
	"os"
	"os/exec"
)

func (p privilege) IsAdmin() bool {
	return os.Geteuid() == 0
}

func (p privilege) Elevate() error {
	// 使用 Polkit 工具 pkexec 来请求授权
	cmd := exec.Command("pkexec", os.Args...)
	// 执行命令
	return cmd.Start()
}
