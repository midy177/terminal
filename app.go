package main

import (
	"context"
	"errors"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"terminal/ent"
	"terminal/ent/folders"
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
			if err != nil {
				log.Printf("error reading from pty: %v", err)
				break
			}
			runtime.EventsEmit(ctx, id, string(buf[:read]))
		}
		runtime.EventsOff(ctx, id)
	}(t, a.ctx, clearFun)
	return nil
}

type HostEntry struct {
	ID       int    `json:"id"`
	IsFolder bool   `json:"is_folder"`
	Label    string `json:"label"`
	Username string `json:"username"`
	Port     uint   `json:"port"`
	Password string `json:"password"`
	FolderID int    `json:"folder_id"`
	KeyID    int    `json:"key_id"`
}

// AddFoldOrHost 添加host
func (a *App) AddFoldOrHost(h *HostEntry) error {
	if h.Port == 0 {
		h.Port = 22
	}
	add := a.db.Hosts.Create().
		SetLabel(h.Label).
		SetUsername(h.Username).
		SetPort(h.Port).
		SetFolderID(h.FolderID).
		SetKeyID(h.KeyID)
	if h.Password != "" {
		add.SetNillablePassword(&h.Password)
	}
	return add.Exec(a.ctx)
}

// DelFoldOrHost 添加DelFoldOrHost
func (a *App) DelFoldOrHost(id int, isFold bool) error {
	return a.db.Hosts.DeleteOneID(id).Exec(a.ctx)
}

// UpdHost 添加host
func (a *App) UpdHost(h *HostEntry) error {
	upd := a.db.Hosts.UpdateOneID(h.ID).
		SetLabel(h.Label).
		SetUsername(h.Username).
		SetPort(h.Port).
		SetFolderID(h.FolderID).
		SetKeyID(h.KeyID)
	if h.Password != "" {
		upd.SetNillablePassword(&h.Password)
	}
	return upd.Exec(a.ctx)
}

// GetHost 通过文件夹ID获取host列表
func (a *App) GetHost(id int) ([]HostEntry, error) {
	all, err := a.db.Folders.GetX(a.ctx, id).QueryHost().All(a.ctx)
	if err != nil {
		return nil, err
	}
	var entries = make([]HostEntry, 0, len(all))
	for _, e := range all {
		entries = append(entries, HostEntry{
			ID:       e.ID,
			Label:    e.Label,
			Username: e.Username,
			Port:     e.Port,
			Password: e.Password,
			FolderID: id,
			KeyID:    e.QueryKey().OnlyIDX(a.ctx),
		})
	}
	return entries, nil
}

// GetFoldsAndHosts 通过文件夹ID获取文件夹列表和主机列表的集合
func (a *App) GetFoldsAndHosts(parentId int) ([]HostEntry, error) {
	allF, err := a.db.Folders.Query().
		Where(folders.ParentID(parentId)).All(a.ctx)
	if err != nil {
		return nil, err
	}

	f, err := a.db.Folders.Get(a.ctx, parentId)
	if err != nil {
		return nil, err
	}
	allH, err := f.QueryHost().All(a.ctx)
	if err != nil {
		return nil, err
	}
	var entries = make([]HostEntry, 0, len(allF)+len(allH))
	for _, e := range allF {
		entries = append(entries, HostEntry{
			ID:       e.ID,
			IsFolder: true,
			Label:    e.Label,
		})
	}
	for _, e := range allH {
		entries = append(entries, HostEntry{
			ID:       e.ID,
			IsFolder: false,
			Label:    e.Label,
			Username: e.Username,
			Port:     e.Port,
			Password: e.Password,
			FolderID: parentId,
			KeyID:    e.QueryKey().OnlyIDX(a.ctx),
		})
	}
	return entries, nil
}
