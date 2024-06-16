package logic

import (
	"errors"
	"fmt"
	"github.com/inhies/go-bytesize"
	"github.com/pkg/sftp"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io"
	"os"
	"path"
	"path/filepath"
)

func (l *Logic) getStpClient(id string) (*sftp.Client, error) {
	t, ok := l.ptyMap.Load(id)
	if !ok {
		return nil, errors.New("pty already released")
	}
	return t.Sftp()
}

type FileInfo struct {
	Name     string // 文件的基本名称
	FullPath string // 完整路径
	Size     string // 常规文件的长度（以字节为单位）;对其他人依赖系统
	Mode     string // 文件模式
	ModTime  int64  // 时间
	IsDir    bool   // abbreviation for Mode().IsDir()
}

func (l *Logic) Dir(id string, dstDir string) ([]FileInfo, error) {
	sftpCli, err := l.getStpClient(id)
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
			FullPath: filepath.Join(dstDir, d.Name()),
			Size:     bytesize.New(float64(d.Size())).String(),
			Mode:     d.Mode().String(),
			ModTime:  d.ModTime().Unix(),
			IsDir:    d.IsDir(),
		})
	}
	return files, nil
}

func (l *Logic) UploadDirectory(id string, dstDir string) error {
	sftpCli, err := l.getStpClient(id)
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
	srcDir, err := runtime.OpenDirectoryDialog(l.Ctx, runtime.OpenDialogOptions{
		Title:                      "选择需要上传的文件夹",
		DefaultDirectory:           "",
		DefaultFilename:            "",
		Filters:                    nil,
		ShowHiddenFiles:            true,
		CanCreateDirectories:       true,
		ResolvesAliases:            false,
		TreatPackagesAsDirectories: false,
	})
	if err != nil {
		return err
	}
	if srcDir == "" {
		return errors.New("src dir is empty")
	}
	return uploadDirectory(sftpCli, srcDir, dstDir)
}

func (l *Logic) UploadMultipleFiles(id string, dstDir string) error {
	sftpCli, err := l.getStpClient(id)
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
	files, err := runtime.OpenMultipleFilesDialog(l.Ctx, runtime.OpenDialogOptions{
		Title: "选择",
	})
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return errors.New("没有选择文件")
	}
	for _, f := range files {
		err := uploadFile(sftpCli, f, dstDir)
		if err != nil {
			_, _ = runtime.MessageDialog(l.Ctx, runtime.MessageDialogOptions{
				Title:   "It's your turn!",
				Message: fmt.Sprintf("上传文件失败: %s", err.Error()),
			})
		}
	}
	return nil
}

func (l *Logic) Download(id string, dst string) error {
	sftpCli, err := l.getStpClient(id)
	if err != nil {
		return err
	}
	dstStat, err := sftpCli.Stat(dst)
	if err != nil {
		return err
	}
	localDir, err := runtime.OpenDirectoryDialog(l.Ctx, runtime.OpenDialogOptions{
		Title: "选择",
	})
	if err != nil {
		return err
	}
	if localDir == "" {
		return errors.New("没有选择本地保存的文件夹")
	}
	if dstStat.IsDir() {
		return downloadDirectory(sftpCli, dst, localDir)
	} else {
		return downloadFile(sftpCli, dst, localDir)
	}
}

func uploadDirectory(sftpClient *sftp.Client, localPath string, remotePath string) error {
	localFiles, err := os.ReadDir(localPath)
	if err != nil {
		return err
	}

	for _, backupDir := range localFiles {
		localFilePath := path.Join(localPath, backupDir.Name())
		remoteFilePath := path.Join(remotePath, backupDir.Name())
		if backupDir.IsDir() {
			if err := sftpClient.Mkdir(remoteFilePath); err != nil {
				return err
			}
			if err := uploadDirectory(sftpClient, localFilePath, remoteFilePath); err != nil {
				return err
			}
		} else {
			if err := uploadFile(sftpClient, path.Join(localPath, backupDir.Name()), remotePath); err != nil {
				return err
			}
		}
	}
	return nil
}

func uploadFile(sftpClient *sftp.Client, localFilePath string, remotePath string) error {
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		return err

	}
	defer srcFile.Close()
	var remoteFileName = path.Base(localFilePath)
	dstFile, err := sftpClient.Create(path.Join(remotePath, remoteFileName))
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

	for _, backupDir := range remoteFiles {
		localFilePath := path.Join(localPath, backupDir.Name())
		remoteFilePath := path.Join(remotePath, backupDir.Name())
		if backupDir.IsDir() {
			if err := os.Mkdir(localFilePath, os.ModeDir); err != nil {
				return err
			}
			if err := downloadDirectory(sftpClient, localFilePath, remoteFilePath); err != nil {
				return err
			}
		} else {
			if err := downloadFile(sftpClient, path.Join(localPath, backupDir.Name()), remotePath); err != nil {
				return err
			}
		}
	}
	return nil
}

func downloadFile(sftpClient *sftp.Client, remoteFilePath, localPath string) error {
	srcFile, err := sftpClient.Open(remoteFilePath)
	if err != nil {
		return err

	}
	defer srcFile.Close()
	var localFileName = path.Base(remoteFilePath)
	dstFile, err := os.Create(path.Join(localPath, localFileName))
	if err != nil {
		return err

	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}