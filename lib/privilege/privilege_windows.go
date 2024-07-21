//go:build windows
// +build windows

package privilege

import (
	"os"
	"os/exec"
	"syscall"
)

// IsAdmin 检查程序是否以管理员权限运行
func (p privilege) IsAdmin() bool {
	f, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	defer f.Close()
	return err == nil
}

func (p privilege) Elevate() error {
	var args []string
	if len(os.Args) > 1 {
		args = os.Args[1:]
	}
	// 准备执行的命令
	cmd := exec.Command(os.Args[0], args...)

	// 使用 runas 参数请求以管理员权限运行
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Token: syscall.Token(0),
	}
	// 执行命令
	return cmd.Start()
}
