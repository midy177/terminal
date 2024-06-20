package logic

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
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
	moveDBFile(sqliteFilePath)
	client, err := ent.Open("sqlite3", fmt.Sprintf("%s%s", sqliteFilePath, "?cache=shared&mode=rwc&_fk=1"))
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

const sqliteFile = "terminal.db"

func getSqliteFilePath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Error getting user config directory:", err)
		return sqliteFile
	}
	dbDir := filepath.Join(configDir, "console.terminal.db")
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		err := os.MkdirAll(dbDir, os.ModePerm)
		if err != nil {
			return sqliteFile
		}
	}
	return filepath.Join(dbDir, sqliteFile)
}

// 迁移文件
func moveDBFile(destFile string) error {
	// 打开源文件
	src, err := os.Open(sqliteFile)
	if err != nil {
		return fmt.Errorf("打开源文件失败: %v", err)
	}
	defer src.Close()

	// 创建目标文件夹（如果不存在）
	destDir := filepath.Dir(destFile)
	err = os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("创建目标文件夹失败: %v", err)
	}

	// 创建目标文件
	dest, err := os.Create(destFile)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %v", err)
	}
	defer dest.Close()

	// 复制文件内容
	_, err = io.Copy(dest, src)
	if err != nil {
		return fmt.Errorf("复制文件内容失败: %v", err)
	}

	// 删除源文件
	err = os.Remove(sqliteFile)
	if err != nil {
		return nil
	}

	return nil
}
