//go:build windows

package shellx

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
)

func getShells() {
	shells = append(shells, SystemShell{
		ID:      "cmd",
		Name:    "CMD",
		Command: "cmd.exe",
		Args:    nil,
		Env:     nil,
		Cwd:     nil,
		Icon:    "/assets/icons/cmd.svg",
	})

	shells = append(shells, SystemShell{
		ID:      "powershell",
		Name:    "PowerShell",
		Command: "powershell.exe",
		Args:    []string{"-NoLogo"},
		Env: []string{
			"TERM=cygwin",
		},
		Cwd:  nil,
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
			fmt.Println("Error getting DistributionName:", err)
			continue
		}

		shell := SystemShell{
			ID:      fmt.Sprintf("wsl-%s", distributionName),
			Name:    fmt.Sprintf("WSL / %s", distributionName),
			Command: "wsl.exe",
			Args:    []string{"-d", distributionName},
			Env: []string{
				"TERM=xterm-256color",
				"COLORTERM=truecolor",
			},
			Cwd:  nil,
			Icon: "/assets/icons/linux.svg",
		}

		shells = append(shells, shell)
	}
	return
}
