package logic

import (
	"context"
	"errors"
	wailsrt "github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"os"
	"runtime"
	"terminal/lib/api"
	"terminal/lib/privilege"
	termx2 "terminal/lib/termx"
)

// GetLocalPtyList 获取本机支持的shell列表
func (l *Logic) GetLocalPtyList() []termx2.SystemShell {
	return termx2.GetShells()
}

// CreateLocalPty 创建本地pty
func (l *Logic) CreateLocalPty(t *termx2.SystemShell) error {
	if _, ok := l.ptyMap.Load(t.ID); ok {
		return errors.New("already exists")
	}
	tPty, err := termx2.NewPTY(t)
	if err != nil {
		return err
	}
	l.ptyMap.Store(t.ID, tPty)
	return l.eventEmitLoop(t.ID)
}

// CreateSshPty 创建ssh pty
func (l *Logic) CreateSshPty(tid string, id, rows, cols int) error {
	one, err := l.db.Hosts.Get(l.Ctx, id)
	if err != nil {
		return err
	}
	var pKey []byte
	if one.KeyID > 0 {
		key, err := one.QueryKey().Only(l.Ctx)
		if err != nil {
			return err
		}
		pKey = key.Content
	}

	term, err := termx2.NewSshPTY(&termx2.SshInfo{
		Username:   one.Username,
		Password:   one.Password,
		Address:    one.Address,
		Port:       one.Port,
		PrivateKey: pKey,
		Height:     rows,
		Width:      cols,
	})
	if err != nil {
		return err
	}
	l.ptyMap.Store(tid, term)
	return l.eventEmitLoop(tid)
}

// CreateSshPtyWithJumper 创建ssh pty with jumper
func (l *Logic) CreateSshPtyWithJumper(id string, tid, jid, rows, cols int) error {
	target, err := l.db.Hosts.Get(l.Ctx, tid)
	if err != nil {
		return err
	}
	var tKey []byte
	if target.KeyID > 0 {
		key, err := target.QueryKey().Only(l.Ctx)
		if err != nil {
			return err
		}
		tKey = key.Content
	}

	jumper, err := l.db.Hosts.Get(l.Ctx, jid)
	if err != nil {
		return err
	}
	var jKey []byte
	if jumper.KeyID > 0 {
		key, err := jumper.QueryKey().Only(l.Ctx)
		if err != nil {
			return err
		}
		jKey = key.Content
	}

	term, err := termx2.NewSshPtyWithJumper(&termx2.SshInfo{
		Username:   target.Username,
		Password:   target.Password,
		Address:    target.Address,
		Port:       target.Port,
		PrivateKey: tKey,
		Height:     rows,
		Width:      cols,
	}, &termx2.SshInfo{
		Username:   jumper.Username,
		Password:   jumper.Password,
		Address:    jumper.Address,
		Port:       jumper.Port,
		PrivateKey: jKey,
		Height:     rows,
		Width:      cols,
	})
	if err != nil {
		return err
	}
	l.ptyMap.Store(id, term)
	return l.eventEmitLoop(id)
}

// ClosePty 关闭pty
func (l *Logic) ClosePty(id string) error {
	l.statMap.Delete(id)
	t, ok := l.ptyMap.LoadAndDelete(id)
	if !ok {
		return errors.New("pty already released")
	}
	return t.(termx2.PtyX).Close()
}

// ResizePty 重置终端大小
func (l *Logic) ResizePty(id string, rows, cols int) error {
	t, ok := l.ptyMap.Load(id)
	if !ok {
		return errors.New("pty already released")
	}
	return t.Resize(rows, cols)
}

// WriteToPty 数据写入pty
func (l *Logic) WriteToPty(id string, data []byte) error {
	t, ok := l.ptyMap.Load(id)
	if !ok {
		return errors.New("pty already released")
	}
	_, err := t.Write(data)
	return err
}

// 推送终端信息到前端
func (l *Logic) eventEmitLoop(id string) error {
	t, ok := l.ptyMap.Load(id)
	if !ok {
		return errors.New("pty already released")
	}
	clearFun := func() {
		_ = t.Close()
		l.ptyMap.Delete(id)
	}
	go func(cPty termx2.PtyX, ctx context.Context, f func()) {
		defer f()
		var buf = make([]byte, 32*1024)
		for {
			read, err := cPty.Read(buf)
			if err != nil {
				log.Printf("error reading from pty: %v\n", err)
				break
			}
			if read > 0 {
				wailsrt.EventsEmit(ctx, id, string(buf[:read]))
			}
		}
		wailsrt.EventsOff(ctx, id)
	}(t, l.Ctx, clearFun)
	return nil
}

func (l *Logic) GetStats(id string) (*api.Stat, error) {
	stat, ok := l.statMap.Load(id)
	if !ok {
		stat = api.NewStats()
		l.statMap.Store(id, stat)
	}
	client, ok := l.ptyMap.Load(id)
	if !ok {
		return nil, errors.New("pty already released")
	}
	sshClient, err := client.Ssh()
	if err != nil {
		return nil, err
	}
	err = stat.GetAllStats(sshClient)
	if err != nil {
		return nil, err
	}
	return stat, nil
}

func (l *Logic) IsRunAsAdmin() bool {
	p := privilege.New()
	return p.IsAdmin()
}

func (l *Logic) RunAsAdmin() error {
	p := privilege.New()
	if p.IsAdmin() {
		return errors.New("already run as admin")
	}
	err := p.Elevate()
	if err != nil {
		return err
	}
	os.Exit(0)
	return nil
}

func (l *Logic) OsGoos() string {
	return runtime.GOOS
}
