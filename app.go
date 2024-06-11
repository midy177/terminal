package main

import (
	"context"
	"errors"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io"
	"log"
	"terminal/ent"
	"terminal/pkg/syncmapx"
	"terminal/termx"
)

// App struct
type App struct {
	ctx    context.Context
	db     *ent.Client
	ptyMap *syncmapx.Map[string, termx.PtyX]
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		ptyMap: syncmapx.New[string, termx.PtyX](),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	client, err := ent.Open("sqlite3", "terminal.db?cache=shared&mode=rwc&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	a.db = client
}

// GetLocalPtyList 获取本机支持的shell列表
func (a *App) GetLocalPtyList() []termx.SystemShell {
	return termx.GetShells()
}

// CreateLocalPty 创建本地pty
func (a *App) CreateLocalPty(t *termx.SystemShell) error {
	if _, ok := a.ptyMap.Load(t.ID); ok {
		return errors.New("already exists")
	}
	tPty, err := termx.NewPTY(t)
	if err != nil {
		return err
	}
	a.ptyMap.Store(t.ID, tPty)
	return a.eventEmitLoop(t.ID)
}

// CreateSshPty 创建ssh pty
func (a *App) CreateSshPty(tid string, id, rows, cols int) error {
	one, err := a.db.Hosts.Get(a.ctx, id)
	if err != nil {
		return err
	}
	var pKey []byte
	if key, err := one.QueryKey().Only(a.ctx); err == nil && key != nil {
		pKey = key.Content
	}

	term, err := termx.NewSshPTY(one.Username,
		one.Password,
		one.Address,
		one.Port,
		pKey,
		rows,
		cols,
	)
	if err != nil {
		return err
	}
	a.ptyMap.Store(tid, term)
	return a.eventEmitLoop(tid)
}

// ClosePty 关闭pty
func (a *App) ClosePty(id string) error {
	t, ok := a.ptyMap.LoadAndDelete(id)
	if !ok {
		return errors.New("pty already released")
	}
	return t.Close()
}

// ResizePty 重置终端大小
func (a *App) ResizePty(id string, rows, cols int) error {
	t, ok := a.ptyMap.Load(id)
	if !ok {
		return errors.New("pty already released")
	}
	return t.Resize(rows, cols)
}

// WriteToPty 数据写入pty
func (a *App) WriteToPty(id string, data []byte) error {
	t, ok := a.ptyMap.Load(id)
	if !ok {
		return errors.New("pty already released")
	}
	_, err := t.Write(data)
	return err
}

// IsWindows 数据写入pty
func (a *App) IsWindows() bool {
	return false
}

// 推送终端信息到前端
func (a *App) eventEmitLoop(id string) error {
	t, ok := a.ptyMap.Load(id)
	if !ok {
		return errors.New("pty already released")
	}
	clearFun := func() {
		_ = t.Close()
		a.ptyMap.Delete(id)
	}
	go func(cPty termx.PtyX, ctx context.Context, f func()) {
		defer f()
		var buf = make([]byte, 32*1024)
		for {
			read, err := cPty.Read(buf)
			if err != nil && err != io.EOF {
				log.Printf("error reading from pty: %v", err)
				break
			}
			runtime.EventsEmit(ctx, id, string(buf[:read]))
		}
		runtime.EventsOff(ctx, id)
	}(t, a.ctx, clearFun)
	return nil
}
