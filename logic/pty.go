package logic

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"terminal/lib/privilege"
	termx2 "terminal/lib/termx"
	"terminal/lib/utils"
	"time"

	"github.com/sqweek/dialog"
	wailsrt "github.com/wailsapp/wails/v2/pkg/runtime"
)

// GetLocalPtyList 获取本机支持的shell列表
func (l *Logic) GetLocalPtyList() []termx2.SystemShell {
	return termx2.GetShells()
}

// CreateLocalPty 创建本地pty
func (l *Logic) CreateLocalPty(t *termx2.SystemShell) error {
	if _, ok := l.Sessions.Load(t.ID); ok {
		return errors.New("已经存在连接")
	}
	tPty, err := termx2.NewPTY(t)
	if err != nil {
		return err
	}
	enabledRec := atomic.Bool{}
	enabledRec.Store(false)
	l.Sessions.Store(t.ID, &Session{
		Pty:        tPty,
		EnabledRec: &enabledRec,
	})
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
	enabledRec := atomic.Bool{}
	enabledRec.Store(false)
	l.Sessions.Store(tid, &Session{
		Pty:        term,
		EnabledRec: &enabledRec,
	})
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
	enabledRec := atomic.Bool{}
	enabledRec.Store(false)
	l.Sessions.Store(id, &Session{
		Pty:        term,
		EnabledRec: &enabledRec,
	})
	return l.eventEmitLoop(id)
}

// ClosePty 关闭pty
func (l *Logic) ClosePty(id string) error {
	l.Sessions.Delete(id)
	t, ok := l.Sessions.LoadAndDelete(id)
	if !ok {
		return errors.New("连接已释放")
	}
	return t.Pty.Close()
}

// ResizePty 重置终端大小
func (l *Logic) ResizePty(id string, rows, cols int) error {
	t, ok := l.Sessions.Load(id)
	if !ok {
		return errors.New("连接已释放")
	}
	if t.Rec != nil {
		_, _ = t.Rec.Resize(rows, cols)
	}
	return t.Pty.Resize(rows, cols)
}

// WriteToPty 数据写入pty
func (l *Logic) WriteToPty(id string, data []byte) error {
	t, ok := l.Sessions.Load(id)
	if !ok {
		return errors.New("连接已释放")
	}
	_, err := t.Pty.Write(data)
	return err
}

// WriteClipboardToPty 数据写入pty
func (l *Logic) WriteClipboardToPty(id string) error {
	clipText, err := wailsrt.ClipboardGetText(l.Ctx)
	if err != nil {
		return err
	}
	clipText = strings.ReplaceAll(clipText, "\r\n", "\r")
	if len(clipText) == 0 {
		return fmt.Errorf("剪贴板里面没有内容")
	}
	t, ok := l.Sessions.Load(id)
	if !ok {
		return errors.New("连接已释放")
	}
	_, err = t.Pty.Write([]byte(clipText))
	return err
}

func (l *Logic) SetClipTextToClipboard(clipText string) error {
	if len(clipText) == 0 {
		return fmt.Errorf("没有选中内容")
	}
	return wailsrt.ClipboardSetText(l.Ctx, clipText)
}

// 推送终端信息到前端
func (l *Logic) eventEmitLoop(id string) error {
	t, ok := l.Sessions.Load(id)
	if !ok {
		return errors.New("连接已释放")
	}

	go func() {
		defer func() {
			_ = t.Pty.Close()
			l.Sessions.Delete(id)
			wailsrt.EventsOff(l.Ctx, id)
			if t.Rec != nil {
				t.Rec.Close()
				t.Rec = nil
			}
		}()

		buf := make([]byte, 32*1024)
		for {
			read, err := t.Pty.Read(buf)
			if err != nil {
				log.Printf("从pty读取数据失败: %v\n", err)
				return
			}
			if read > 0 {
				data := buf[:read]
				wailsrt.EventsEmit(l.Ctx, id, string(data))
				if t.EnabledRec.Load() {
					_, _ = t.Rec.Write(data)
				}
			}
		}
	}()

	return nil
}

func (l *Logic) StartRec(id string, rows, cols int) (string, error) {
	sess, ok := l.Sessions.Load(id)
	if !ok {
		return "", errors.New("没有创建录屏或连接已经释放")
	}
	if sess.Rec != nil {
		return "", errors.New("录屏已经在进行中")
	}
	if sess.EnabledRec.Load() {
		return "", errors.New("已经开启录屏")
	}

	// 打开文件夹选择对话框
	folderPath, err := dialog.Directory().Title("选择录屏文件存档文件夹").Browse()
	if err != nil {
		if errors.Is(err, dialog.ErrCancelled) {
			return "", errors.New("用户取消了选择")
		}
		return "", errors.New("打开文件夹对话框出错: " + err.Error())
	}

	// terminal_recording
	filename := filepath.Join(folderPath, id+"_"+time.Now().Format("20060102150405")+".cast")

	sess.Rec, err = utils.NewRecorder(filename, rows, cols)
	if err != nil {
		return "", err
	}
	sess.EnabledRec.Store(true)
	return filename, nil
}

func (l *Logic) StopRec(id string) error {
	sess, ok := l.Sessions.Load(id)
	if !ok {
		return errors.New("连接已关闭")
	}
	if sess.EnabledRec.Load() {
		sess.EnabledRec.Store(false)
	}
	if sess.Rec == nil {
		return errors.New("没有创建录屏")
	}
	sess.Rec.Close()
	sess.Rec = nil
	return nil
}

//func (l *Logic) GetStats(id string) (*api.Stat, error) {
//	sess, ok := l.Sessions.Load(id)
//	if ok && sess.stat != nil { {
//		stat = api.NewStats()
//		l.Sessions.Store(id, stat)
//	}
//	client, ok := l.ptyMap.Load(id)
//	if !ok {
//		return nil, errors.New("pty already released")
//	}
//	sshClient, err := client.Ssh()
//	if err != nil {
//		return nil, err
//	}
//	err = stat.GetAllStats(sshClient)
//	if err != nil {
//		return nil, err
//	}
//	return stat, nil
//}

func (l *Logic) IsRunAsAdmin() bool {
	p := privilege.New()
	return p.IsAdmin()
}

func (l *Logic) RunAsAdmin() error {
	p := privilege.New()
	if p.IsAdmin() {
		return errors.New("已经是管理员")
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
