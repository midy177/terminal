package logic

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"terminal/ent"
	"terminal/pkg/syncmapx"
	"terminal/termx"
)

// Logic struct
type Logic struct {
	Ctx    context.Context
	db     *ent.Client
	ptyMap *syncmapx.Map[string, termx.PtyX]
}

// NewApp creates a new App application struct
func NewApp() *Logic {
	l := &Logic{
		ptyMap: syncmapx.New[string, termx.PtyX](),
	}
	sqliteFilePath := getSqliteFilePath()
	client, err := ent.Open("sqlite3", sqliteFilePath)
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	l.db = client
	return l
}

const sqliteFile = "terminal.db?cache=shared&mode=rwc&_fk=1"

func getSqliteFilePath() string {
	if isGoRun() {
		return sqliteFile
	}
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Error getting user config directory:", err)
		return sqliteFile
	}
	return filepath.Join(configDir, sqliteFile)
}

func isGoRun() bool {
	// 获取可执行文件的路径
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path:", err)
		return false
	}
	// 获取可执行文件的目录
	exeDir := filepath.Dir(exePath)

	// 检查目录是否包含临时目录路径的一部分
	tempDir := os.TempDir()
	return strings.HasPrefix(exeDir, tempDir)
}
