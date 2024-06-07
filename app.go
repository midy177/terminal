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
	defer client.Close()
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
func (a *App) ClosePty(id string) error {
	t, ok := a.ptyMap.LoadAndDelete(id)
	if ok {
		return errors.New("already exists")
	}
	return t.Close()
}

func (a *App) ResizePty(id string, rows, cols int) error {
	t, ok := a.ptyMap.LoadAndDelete(id)
	if ok {
		return errors.New("already exists")
	}
	return t.Resize(rows, cols)
}

func (a *App) CreateSshPty(id int) error {

	return nil
}

func (a *App) eventEmitLoop(id string) error {
	t, ok := a.ptyMap.Load(id)
	if ok {
		return errors.New("already exists")
	}
	go func(cPty termx.PtyX, ctx context.Context) {
		defer cPty.Close()
		var buf = make([]byte, 32*1024)
		for {
			read, err := cPty.Read(buf)
			if err != nil && err != io.EOF {
				break
			}
			runtime.EventsEmit(ctx, id, buf[:read])
		}
	}(t, a.ctx)
	return nil
}
