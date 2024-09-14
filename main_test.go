package main

import (
	"context"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strings"
	"terminal/ent"
	"terminal/lib/utils"
	"testing"
)

func TestName(t *testing.T) {
	ctx := context.Background()
	client, err := ent.Open("sqlite3", "terminal.db?cache=shared&mode=rwc&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	err = client.Folders.Create().SetLabel("root").Exec(ctx)
	if err != nil {
		log.Fatalf("failed creating folder: %v", err)
	}
	all, err := client.Folders.Query().All(ctx)
	if err != nil {
		log.Fatalf("failed querying folders: %v", err)
	}
	for _, v := range all {
		fmt.Println(v)
	}
}

func TestLogger(t *testing.T) {
	logFilter, err := utils.NewLogFilter("ssh_session.log")
	if err != nil {
		log.Fatalf("无法创建日志过滤器: %v", err)
	}
	defer logFilter.Close()

	// 这里应该是你的 SSH 会话输出
	input := `
常规输出
rz waiting to receive.
文件传输数据...
文件传输数据...
50% complete
文件传输数据...
100% complete
Transfer complete
更多常规输出
`
	logFilter.ProcessOutput(strings.NewReader(input))
}
