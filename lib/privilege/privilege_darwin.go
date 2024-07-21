//go:build darwin
// +build darwin

package privilege

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func (p privilege) IsAdmin() bool {
	return os.Geteuid() == 0
}

func (p privilege) Elevate() error {
	// 检查程序是否以管理员权限运行
	if os.Geteuid() == 0 {
		return nil
	}
	// 组装 AppleScript
	script := fmt.Sprintf(`do shell script "sudo %s" with administrator privileges`,
		strings.Join(os.Args[0:], " "))

	// 构建执行 AppleScript 的命令

	cmd := exec.Command("osascript", "-e", script)
	return cmd.Start()
}
