package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"terminal/termx"
)

// GetLocalPtyList 获取本机支持的shell列表
func (l *Logic) GetLocalPtyList() []termx.SystemShell {
	return termx.GetShells()
}

// CreateLocalPty 创建本地pty
func (l *Logic) CreateLocalPty(t *termx.SystemShell) error {
	if _, ok := l.ptyMap.Load(t.ID); ok {
		return errors.New("already exists")
	}
	tPty, err := termx.NewPTY(t)
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
		fmt.Println(one.KeyID)
		key, err := one.QueryKey().Only(l.Ctx)
		if err == nil {
			return err
		}
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
	l.ptyMap.Store(tid, term)
	return l.eventEmitLoop(tid)
}

// ClosePty 关闭pty
func (l *Logic) ClosePty(id string) error {
	t, ok := l.ptyMap.LoadAndDelete(id)
	if !ok {
		return errors.New("pty already released")
	}
	return t.Close()
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
	}(t, l.Ctx, clearFun)
	return nil
}
