package termx

import (
	"fmt"
	"io"
	"os"
	"sync/atomic"
	"terminal/lib/trzsz"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

const (
	defaultTermType = "xterm-256color"
	sshDialTimeout  = 15 * time.Second
	ttySpeed        = 14400
)

type sshSession struct {
	client      *ssh.Client
	session     *ssh.Session
	stdin       io.WriteCloser
	stdout      io.Reader
	stderr      io.Reader
	sftp        *sftp.Client
	trzszFilter *trzsz.TrzszFilter
	closed      *atomic.Bool
	clear       func()
}

func (s *sshSession) Ssh() (*ssh.Client, error) {
	return s.client, nil
}

// Sftp create sftp client
func (s *sshSession) Sftp() (*sftp.Client, error) {
	if s.sftp != nil {
		return s.sftp, nil
	}
	var err error
	if s.sftp, err = sftp.NewClient(s.client); err != nil {
		return nil, err
	}
	return s.sftp, nil
}

// CloseSftp close sftp client
func (s *sshSession) CloseSftp() error {
	if s.sftp != nil {
		defer func() {
			s.sftp = nil
		}()
		return s.sftp.Close()
	}
	return nil
}

func (s *sshSession) Resize(rows, cols int) error {
	s.trzszFilter.SetTerminalColumns(int32(cols))
	return s.session.WindowChange(rows, cols)
}

func (s *sshSession) Read(p []byte) (n int, err error) {
	if s.closed.Load() {
		fmt.Printf("SSH 连接已关闭: %s\n", err)
		return 0, io.EOF
	}
	or, err := s.stdout.Read(p)
	if err != nil {
		return 0, err
	}
	if or > 0 {
		return or, nil
	}
	return s.stderr.Read(p)
}

func (s *sshSession) Write(p []byte) (n int, err error) {
	if s.closed.Load() {
		fmt.Printf("SSH 连接已关闭: %s\n", err)
		return 0, io.EOF
	}
	return s.stdin.Write(p)
}

func (s *sshSession) Close() error {
	if s.closed.CompareAndSwap(false, true) {
		s.clear()
	}
	return nil
}

type SshInfo struct {
	Username, Password, Address string
	Port                        uint
	PrivateKey                  []byte
	Height, Width               int
}

// 辅助函数，用于关闭所有资源
func closeAll(closers ...io.Closer) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	for _, closer := range closers {
		if closer != nil {
			_ = closer.Close()
		}
	}
}

// 创建SSH会话的通用函数
func createSSHSession(client *ssh.Client, targetInfo *SshInfo) (*sshSession, error) {
	termType := os.Getenv("TERM")
	if termType == "" {
		termType = defaultTermType
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("无法创建SSH会话: %w", err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.ECHOCTL:       0,
		ssh.TTY_OP_ISPEED: ttySpeed,
		ssh.TTY_OP_OSPEED: ttySpeed,
	}
	if err := session.RequestPty(termType, targetInfo.Height, targetInfo.Width, modes); err != nil {
		session.Close()
		return nil, fmt.Errorf("无法请求伪终端: %w", err)
	}

	serverIn, err := session.StdinPipe()
	if err != nil {
		session.Close()
		return nil, fmt.Errorf("无法获取标准输入流: %w", err)
	}
	serverOut, err := session.StdoutPipe()
	if err != nil {
		session.Close()
		return nil, fmt.Errorf("无法获取标准输出流: %w", err)
	}
	stderr, err := session.StderrPipe()
	if err != nil {
		session.Close()
		return nil, fmt.Errorf("无法获取标准错误流: %w", err)
	}

	clientIn, stdinPipe := io.Pipe()
	stdoutPipe, clientOut := io.Pipe()
	trzszFilter := trzsz.NewTrzszFilter(clientIn, clientOut, serverIn, serverOut,
		trzsz.TrzszOptions{
			TerminalColumns: int32(targetInfo.Width),
			EnableZmodem:    true,
		})

	if err := session.Shell(); err != nil {
		session.Close()
		return nil, fmt.Errorf("无法启动SSH shell: %w", err)
	}

	closed := &atomic.Bool{}
	go func() {
		err := session.Wait()
		if err != nil {
			fmt.Printf("SSH shell 退出: %s\n", err)
		} else {
			fmt.Printf("SSH shell 正常退出\n")
		}

		if closed.CompareAndSwap(false, true) {
			closeAll(session, clientIn, stdinPipe, stdoutPipe, clientOut, serverIn, client)
		}
	}()

	return &sshSession{
		client:      client,
		session:     session,
		stdin:       stdinPipe,
		stdout:      stdoutPipe,
		stderr:      stderr,
		trzszFilter: trzszFilter,
		clear: func() {
			closeAll(session, clientIn, stdinPipe, stdoutPipe, clientOut, serverIn, client)
		},
		closed: closed,
	}, nil
}

func NewSshPTY(targetInfo *SshInfo) (PtyX, error) {
	sshConfig := &ssh.ClientConfig{
		User:            targetInfo.Username,
		Auth:            []ssh.AuthMethod{ssh.Password(targetInfo.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         sshDialTimeout,
	}

	if targetInfo.PrivateKey != nil {
		signer, err := ssh.ParsePrivateKey(targetInfo.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("解析私钥失败: %w", err)
		}
		sshConfig.Auth = append(sshConfig.Auth, ssh.PublicKeys(signer))
	}

	sshClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", targetInfo.Address, targetInfo.Port), sshConfig)
	if err != nil {
		return nil, fmt.Errorf("无法创建SSH连接: %w", err)
	}

	return createSSHSession(sshClient, targetInfo)
}

func NewSshPtyWithJumper(targetInfo, jumpInfo *SshInfo) (PtyX, error) {
	jumpConfig := &ssh.ClientConfig{
		User:            jumpInfo.Username,
		Auth:            []ssh.AuthMethod{ssh.Password(jumpInfo.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         sshDialTimeout,
	}

	if jumpInfo.PrivateKey != nil {
		signer, err := ssh.ParsePrivateKey(jumpInfo.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("解析跳板机私钥失败: %w", err)
		}
		jumpConfig.Auth = append(jumpConfig.Auth, ssh.PublicKeys(signer))
	}

	jumpConn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", jumpInfo.Address, jumpInfo.Port), jumpConfig)
	if err != nil {
		return nil, fmt.Errorf("无法连接到跳板机: %w", err)
	}

	targetConfig := &ssh.ClientConfig{
		User:            targetInfo.Username,
		Auth:            []ssh.AuthMethod{ssh.Password(targetInfo.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         sshDialTimeout,
	}

	if targetInfo.PrivateKey != nil {
		signer, err := ssh.ParsePrivateKey(targetInfo.PrivateKey)
		if err != nil {
			jumpConn.Close()
			return nil, fmt.Errorf("解析目标主机私钥失败: %w", err)
		}
		targetConfig.Auth = append(targetConfig.Auth, ssh.PublicKeys(signer))
	}

	targetAddr := fmt.Sprintf("%s:%d", targetInfo.Address, targetInfo.Port)
	targetConn, err := jumpConn.Dial("tcp", targetAddr)
	if err != nil {
		jumpConn.Close()
		return nil, fmt.Errorf("无法通过跳板机连接到目标主机: %w", err)
	}

	ncc, chans, reqs, err := ssh.NewClientConn(targetConn, targetAddr, targetConfig)
	if err != nil {
		targetConn.Close()
		jumpConn.Close()
		return nil, fmt.Errorf("无法创建客户端连接: %w", err)
	}

	targetClient := ssh.NewClient(ncc, chans, reqs)

	return createSSHSession(targetClient, targetInfo)
}
