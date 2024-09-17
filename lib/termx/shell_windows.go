//go:build windows

package termx

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
)

func getShellsOld() {
	home := startDir()
	shells = append(shells, SystemShell{
		//ID:      "cmd",
		Name:    "CMD",
		Command: "cmd.exe",
		Args:    nil,
		Env:     nil,
		Cwd:     home,
		Icon:    "/assets/icons/cmd.svg",
	})

	shells = append(shells, SystemShell{
		//ID:      "powershell",
		Name:    "PowerShell",
		Command: "powershell.exe",
		Args:    []string{"-NoLogo"},
		Env: []string{
			"TERM=cygwin",
		},
		Cwd:  home,
		Icon: "/assets/icons/powershell.svg",
	})

	lxssPath := `Software\Microsoft\Windows\CurrentVersion\Lxss`
	k, err := registry.OpenKey(registry.CURRENT_USER, lxssPath, registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		fmt.Println("Error opening registry key:", err)
		return
	}
	defer k.Close()

	keys, err := k.ReadSubKeyNames(-1)
	if err != nil {
		fmt.Println("Error reading subkeys:", err)
		return
	}

	for _, key := range keys {
		subKeyPath := fmt.Sprintf(`%s\%s`, lxssPath, key)
		subKey, err := registry.OpenKey(registry.CURRENT_USER, subKeyPath, registry.READ)
		if err != nil {
			fmt.Println("Error opening subkey:", err)
			continue
		}

		distributionName, _, err := subKey.GetStringValue("DistributionName")
		subKey.Close()
		if err != nil {
			continue
		}

		shell := SystemShell{
			//ID:      fmt.Sprintf("wsl-%s", distributionName),
			Name:    fmt.Sprintf("WSL(%s)", distributionName),
			Command: "wsl.exe",
			Args:    []string{"-d", distributionName},
			Env: []string{
				"TERM=xterm-256color",
				"COLORTERM=truecolor",
			},
			Cwd:  home,
			Icon: "/assets/icons/linux.svg",
		}

		shells = append(shells, shell)
	}
	return
}

func getShells() {
	term := os.Getenv("TERM")
	if term == "" {
		term = "xterm-256color"
	}
	home := startDir()
	user := os.Getenv("USERNAME")
	hostname := getShortHostname()
	title := user
	if hostname != "" {
		title += "@" + hostname
	}

	// 添加 CMD
	shells = append(shells, SystemShell{
		Name:    "CMD",
		Command: "cmd.exe",
		Args:    nil,
		Env:     []string{"TERM=" + term},
		Cwd:     home,
		Icon:    "/assets/icons/windows.svg",
	})

	// 添加 PowerShell
	powershellPath, err := exec.LookPath("powershell.exe")
	if err == nil {
		shells = append(shells, SystemShell{
			Name:    "PowerShell",
			Command: powershellPath,
			Args:    nil,
			Env:     []string{"TERM=" + term},
			Cwd:     home,
			Icon:    "/assets/icons/powershell.svg",
		})
	}

	// 检查是否安装了 Git Bash
	gitBashPath := "C:\\Program Files\\Git\\bin\\bash.exe"
	if _, err := os.Stat(gitBashPath); err == nil {
		shells = append(shells, SystemShell{
			Name:    "Git Bash",
			Command: gitBashPath,
			Args:    nil,
			Env:     []string{"TERM=" + term},
			Cwd:     home,
			Icon:    "/assets/icons/git-bash.svg",
		})
	}

	lxssPath := `Software\Microsoft\Windows\CurrentVersion\Lxss`
	k, err := registry.OpenKey(registry.CURRENT_USER, lxssPath, registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		fmt.Println("Error opening registry key:", err)
		return
	}
	defer k.Close()

	keys, err := k.ReadSubKeyNames(-1)
	if err != nil {
		fmt.Println("Error reading subkeys:", err)
		return
	}

	for _, key := range keys {
		subKeyPath := fmt.Sprintf(`%s\%s`, lxssPath, key)
		subKey, err := registry.OpenKey(registry.CURRENT_USER, subKeyPath, registry.READ)
		if err != nil {
			fmt.Println("Error opening subkey:", err)
			continue
		}

		distributionName, _, err := subKey.GetStringValue("DistributionName")
		subKey.Close()
		if err != nil {
			continue
		}

		shell := SystemShell{
			//ID:      fmt.Sprintf("wsl-%s", distributionName),
			Name:    fmt.Sprintf("WSL(%s)", distributionName),
			Command: "wsl.exe",
			Args:    []string{"-d", distributionName},
			Env: []string{
				"TERM=xterm-256color",
				"COLORTERM=truecolor",
			},
			Cwd:  home,
			Icon: "/assets/icons/linux.svg",
		}

		shells = append(shells, shell)
	}
}

func getShortHostname() string {
	cmd := exec.Command("hostname")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}
