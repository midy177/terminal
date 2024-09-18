package logic

import (
	"errors"
	"github.com/inhies/go-bytesize"
	"github.com/pkg/sftp"
	"github.com/sqweek/dialog"
	"io"
	"os"
	"path"
	"path/filepath"
)

// getSftpClient 获取sftp实例
func (l *Logic) getSftpClient(id string) (*sftp.Client, error) {
	t, ok := l.Sessions.Load(id)
	if !ok {
		return nil, errors.New("pty already released")
	}
	return t.Pty.Sftp()
}

// CloseSftpClient 关闭sftp实例
func (l *Logic) CloseSftpClient(id string) error {
	t, ok := l.Sessions.Load(id)
	if !ok {
		return errors.New("pty already released")
	}
	return t.Pty.CloseSftp()
}

type FileInfo struct {
	Name     string `json:"name"`      // 文件的基本名称
	FullPath string `json:"full_path"` // 完整路径
	Size     string `json:"size"`      // 常规文件的长度（以字节为单位）;对其他人依赖系统
	Mode     string `json:"mode"`      // 文件模式
	ModTime  int64  `json:"mod_time"`  // 时间
	IsDir    bool   `json:"is_dir"`    // abbreviation for Mode().IsDir()
}

// SftpHomeDir 获取sftp Home路径
func (l *Logic) SftpHomeDir(id string) (string, error) {
	sftpCli, err := l.getSftpClient(id)
	if err != nil {
		return "", err
	}
	return sftpCli.Getwd()
}

// SftpDir sftp获取文件夹列表
func (l *Logic) SftpDir(id string, dstDir string) ([]FileInfo, error) {
	sftpCli, err := l.getSftpClient(id)
	if err != nil {
		return nil, err
	}
	if dstDir == "" {
		wd, err := sftpCli.Getwd()
		if err != nil {
			return nil, err
		}
		dstDir = wd
	}
	dirs, err := sftpCli.ReadDir(dstDir)
	if err != nil {
		return nil, err
	}
	var files = make([]FileInfo, 0, len(dirs))
	for _, d := range dirs {
		files = append(files, FileInfo{
			Name:     d.Name(),
			FullPath: path.Join(dstDir, d.Name()),
			Size:     bytesize.New(float64(d.Size())).String(),
			Mode:     d.Mode().String(),
			ModTime:  d.ModTime().Unix(),
			IsDir:    d.IsDir(),
		})
	}
	return files, nil
}

// SftpUploadDirectory sftp上传文件夹
func (l *Logic) SftpUploadDirectory(id string, dstDir string) error {
	sftpCli, err := l.getSftpClient(id)
	if err != nil {
		return err
	}
	dstDirStat, err := sftpCli.Stat(dstDir)
	if err != nil {
		return err
	}
	if !dstDirStat.IsDir() {
		return errors.New("dst dir is not l directory")
	}
	// 打开文件夹选择对话框
	srcDir, err := dialog.Directory().Title("选择要上传的文件夹").Browse()
	if err != nil {
		if errors.Is(err, dialog.ErrCancelled) {
			return errors.New("用户取消了选择")
		}
		return errors.New("打开文件夹对话框出错: " + err.Error())
	}
	if err != nil {
		return err
	}
	if srcDir == "" {
		return errors.New("src dir is empty")
	}
	return uploadDirectory(sftpCli, srcDir, dstDir)
}

// SftpUploadMultipleFiles sftp上传多文件
func (l *Logic) SftpUploadMultipleFiles(id string, dstDir string) error {
	sftpCli, err := l.getSftpClient(id)
	if err != nil {
		return err
	}
	dstDirStat, err := sftpCli.Stat(dstDir)
	if err != nil {
		return err
	}
	if !dstDirStat.IsDir() {
		return errors.New("dst dir is not l directory")
	}
	// 打开文件夹选择对话框
	file, err := dialog.File().Title("选择上传文件").Load()
	if err != nil {
		if errors.Is(err, dialog.ErrCancelled) {
			return errors.New("用户取消了选择")
		}
		return errors.New("打开文件夹对话框出错: " + err.Error())
	}
	fname := filepath.Base(file)
	remoteFilePath := path.Join(dstDir, fname)
	err = uploadFile(sftpCli, file, remoteFilePath)
	if err != nil {
		return err
	}
	return nil
}

// SftpDownload sftp下载文件/文件夹
func (l *Logic) SftpDownload(id string, dst string) error {
	sftpCli, err := l.getSftpClient(id)
	if err != nil {
		return err
	}
	dstStat, err := sftpCli.Stat(dst)
	if err != nil {
		return err
	}
	// 打开文件夹选择对话框
	localDir, err := dialog.Directory().Title("选择保存文件夹").Browse()
	if err != nil {
		if errors.Is(err, dialog.ErrCancelled) {
			return errors.New("用户取消了选择")
		}
		return errors.New("打开文件夹对话框出错: " + err.Error())
	}
	if localDir == "" {
		return errors.New("没有选择本地保存的文件夹")
	}
	if dstStat.IsDir() {
		return downloadDirectory(sftpCli, dst, localDir)
	} else {
		localFilePath := filepath.Join(localDir, path.Base(dst))
		return downloadFile(sftpCli, dst, localFilePath)
	}
}

// SftpDelete sftp删除文件/文件夹
func (l *Logic) SftpDelete(id string, dst string) error {
	sftpCli, err := l.getSftpClient(id)
	if err != nil {
		return err
	}
	return sftpCli.Remove(dst)
}

func uploadDirectory(sftpClient *sftp.Client, localPath string, remotePath string) error {
	localFiles, err := os.ReadDir(localPath)
	if err != nil {
		return err
	}
	for _, backupDir := range localFiles {
		localFilePath := filepath.Join(localPath, backupDir.Name())
		remoteFilePath := path.Join(remotePath, backupDir.Name())
		if backupDir.IsDir() {
			if err := sftpClient.Mkdir(remoteFilePath); err != nil {
				return err
			}
			if err := uploadDirectory(sftpClient, localFilePath, remoteFilePath); err != nil {
				return err
			}
		} else {
			if err := uploadFile(sftpClient, localFilePath, remoteFilePath); err != nil {
				return err
			}
		}
	}
	return nil
}

func uploadFile(sftpClient *sftp.Client, localFilePath string, remoteFilePath string) error {
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		return err

	}
	defer srcFile.Close()
	dstFile, err := sftpClient.Create(remoteFilePath)
	if err != nil {
		return err

	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func downloadDirectory(sftpClient *sftp.Client, remotePath, localPath string) error {
	remoteFiles, err := sftpClient.ReadDir(remotePath)
	if err != nil {
		return err
	}
	folderName := filepath.Base(filepath.Clean(remotePath))
	localPath = filepath.Join(localPath, folderName)
	if err := os.Mkdir(localPath, os.ModeDir); err != nil {
		return err
	}
	for _, backupDir := range remoteFiles {
		localFilePath := filepath.Join(localPath, backupDir.Name())
		remoteFilePath := path.Join(remotePath, backupDir.Name())
		if backupDir.IsDir() {
			if err := os.Mkdir(localFilePath, os.ModeDir); err != nil {
				return err
			}
			if err := downloadDirectory(sftpClient, remoteFilePath, localFilePath); err != nil {
				return err
			}
		} else {
			if err := downloadFile(sftpClient, remoteFilePath, localFilePath); err != nil {
				return err
			}
		}
	}
	return nil
}

func downloadFile(sftpClient *sftp.Client, remoteFilePath, localFilePath string) error {
	srcFile, err := sftpClient.Open(remoteFilePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := os.Create(localFilePath)
	if err != nil {
		return err

	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func calculateLocalTotalSize(localDir ...string) int64 {
	var totalSize int64
	for _, dir := range localDir {
		_ = filepath.Walk(dir, func(localPath string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() {
				totalSize += info.Size()
			}
			return nil
		})
	}
	return totalSize
}
func calculateRemoteTotalSize(client *sftp.Client, remoteDir string) int64 {
	var totalSize int64
	walker := client.Walk(remoteDir)
	for walker.Step() {
		if err := walker.Err(); err == nil && !walker.Stat().IsDir() {
			totalSize += walker.Stat().Size()
		}
	}
	return totalSize
}

//type progressEvent struct {
//	IsDownload bool
//}
//
//func (r *progressEvent) Write(p []byte) (n int, err error) {
//	//TODO implement me
//	panic("implement me")
//}
//func (r *progressEvent) Finish() {
//
//}
