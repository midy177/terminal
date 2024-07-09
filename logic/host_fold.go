package logic

import (
	"errors"
	"terminal/ent/folders"
	"terminal/ent/hosts"
	"terminal/ent/keys"
	"terminal/ent/predicate"
)

type HostEntry struct {
	ID       int    `json:"id"`
	IsFolder bool   `json:"is_folder"`
	Label    string `json:"label"`
	Username string `json:"username"`
	Address  string `json:"address"`
	Port     uint   `json:"port"`
	Password string `json:"password"`
	FolderID int    `json:"folder_id"`
	KeyID    int    `json:"key_id"`
}

// AddFoldOrHost 添加host
func (l *Logic) AddFoldOrHost(h *HostEntry) error {
	if h.IsFolder {
		add := l.db.Folders.
			Create().
			SetLabel(h.Label)
		if h.FolderID > 0 {
			add.SetParentID(h.FolderID)
		}
		return add.Exec(l.Ctx)
	} else {
		if h.Port == 0 {
			h.Port = 22
		}
		add := l.db.Hosts.Create().
			SetLabel(h.Label).
			SetUsername(h.Username).
			SetAddress(h.Address).
			SetPort(h.Port)
		if h.FolderID != 0 {
			add.SetFolderID(h.FolderID)
		}
		if h.Password != "" {
			add.SetNillablePassword(&h.Password)
		}
		if h.KeyID != 0 {
			add.SetKeyID(h.KeyID)
		}
		return add.Exec(l.Ctx)
	}
}

// DelFoldOrHost 添加DelFoldOrHost
func (l *Logic) DelFoldOrHost(id int, isFold bool) error {
	if isFold {
		return l.db.Folders.DeleteOneID(id).Exec(l.Ctx)
	} else {
		return l.db.Hosts.DeleteOneID(id).Exec(l.Ctx)
	}
}

// UpdFoldOrHost 添加host或目录
func (l *Logic) UpdFoldOrHost(h *HostEntry) error {
	if h.IsFolder {
		upd := l.db.Folders.
			UpdateOneID(h.ID).
			SetLabel(h.Label)
		if h.FolderID == 0 {
			upd.ClearParentID()
		} else {
			upd.SetParentID(h.FolderID)
		}
		return upd.Exec(l.Ctx)
	} else {
		upd := l.db.Hosts.UpdateOneID(h.ID).
			SetLabel(h.Label).
			SetUsername(h.Username).
			SetAddress(h.Address).
			SetPort(h.Port).
			SetPassword(h.Password)
		if h.FolderID == 0 {
			upd.ClearFolderID()
		} else {
			upd.SetFolderID(h.FolderID)
		}
		if h.KeyID == 0 {
			upd.ClearKeyID()
		} else {
			upd.SetKeyID(h.KeyID)
		}
		return upd.Exec(l.Ctx)
	}
}

// GetFolds 通过文件夹ID获取host列表
func (l *Logic) GetFolds() ([]HostEntry, error) {
	all, err := l.db.Folders.Query().All(l.Ctx)
	if err != nil {
		return nil, err
	}
	var entries = make([]HostEntry, 0, len(all)+1)
	entries = append(entries, HostEntry{
		ID:    0,
		Label: "根",
	})
	for _, e := range all {
		entries = append(entries, HostEntry{
			ID:       e.ID,
			Label:    e.Label,
			FolderID: e.ParentID,
		})
	}
	return entries, nil
}

// GetHosts 通过文件夹ID获取host列表
func (l *Logic) GetHosts(id int) ([]HostEntry, error) {
	all, err := l.db.Folders.GetX(l.Ctx, id).QueryHost().All(l.Ctx)
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
			KeyID:    e.QueryKey().OnlyIDX(l.Ctx),
		})
	}
	return entries, nil
}

// GetFoldsAndHosts 通过文件夹ID获取文件夹列表和主机列表的集合
func (l *Logic) GetFoldsAndHosts(parentId int) ([]HostEntry, error) {
	var (
		pf predicate.Folders
		hf predicate.Hosts
	)
	if parentId > 0 {
		pf = folders.ParentID(parentId)
		hf = hosts.FolderID(parentId)
	} else {
		pf = folders.ParentIDIsNil()
		hf = hosts.FolderIDIsNil()
	}
	// 获取文件夹
	allF, err := l.db.Folders.Query().
		Where(pf).All(l.Ctx)
	if err != nil {
		return nil, err
	}
	allH, err := l.db.Hosts.Query().Where(hf).All(l.Ctx)
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
			Address:  e.Address,
			Port:     e.Port,
			Password: e.Password,
			FolderID: parentId,
			KeyID:    e.KeyID,
		})
	}
	return entries, nil
}

type KeyEntry struct {
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// 标记
	Label string `json:"label,omitempty"`
	// 私钥信息
	Content string `json:"content,omitempty"`
}

// GetKeyList 获取私钥列表
func (l *Logic) GetKeyList(withContent bool) ([]KeyEntry, error) {
	all, err := l.db.Keys.Query().All(l.Ctx)
	if err != nil {
		return nil, err
	}
	var entries = make([]KeyEntry, 0, len(all)+1)
	entries = append(entries, KeyEntry{
		ID:    0,
		Label: "空",
	})
	for _, e := range all {
		entry := KeyEntry{
			ID:    e.ID,
			Label: e.Label,
		}
		if withContent {
			entry.Content = string(e.Content)
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

// AddKey 添加私钥
func (l *Logic) AddKey(k *KeyEntry) error {
	count, err := l.db.Keys.Query().Where(keys.ContentEQ([]byte(k.Content))).Count(l.Ctx)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("key already exists")
	}
	return l.db.Keys.
		Create().
		SetLabel(k.Label).
		SetContent([]byte(k.Content)).
		Exec(l.Ctx)
}

// UpdKey 修改私钥
func (l *Logic) UpdKey(k *KeyEntry) error {
	return l.db.Keys.
		UpdateOneID(k.ID).
		SetLabel(k.Label).
		SetContent([]byte(k.Content)).
		Exec(l.Ctx)
}

// DelKey 删除私钥
func (l *Logic) DelKey(id int) error {
	return l.db.Keys.
		DeleteOneID(id).
		Exec(l.Ctx)
}
