package termx

import (
	"fmt"
	"github.com/pkg/sftp"
	"github.com/trzsz/trzsz-go/trzsz"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"time"
)

type sshSession struct {
	client      *ssh.Client
	session     *ssh.Session
	stdin       io.WriteCloser
	stdout      io.Reader
	stderr      io.Reader
	sftp        *sftp.Client
	trzszFilter *trzsz.TrzszFilter
	clear       func()
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

func (s *sshSession) Resize(rows, cols int) error {
	s.trzszFilter.SetTerminalColumns(int32(cols))
	return s.session.WindowChange(rows, cols)
}

func (s *sshSession) Read(p []byte) (n int, err error) {
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
	return s.stdin.Write(p)
}

func (s *sshSession) Close() error {
	s.clear()
	return nil
}

func NewSshPTY(username, password, address string, port uint, privateKey []byte, height, width int) (PtyX, error) {
	termType := os.Getenv("TERM")
	if termType == "" {
		termType = "xterm-256color"
	}
	sshConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		//Specify the host key verification callback function
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 15,
	}
	if privateKey != nil {
		// special case and we got a key
		signer, err := ssh.ParsePrivateKey(privateKey)
		if err != nil {
			return nil, err
		}
		sshConfig.Auth = append(sshConfig.Auth, ssh.PublicKeys(signer))
	}

	sshClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", address, port), sshConfig)
	if err != nil {
		fmt.Printf("%s%s\n\r", "Unable to create SSH connection: ", err)
		return nil, err
	}

	session, err := sshClient.NewSession()
	if err != nil {
		fmt.Printf("Unable to create SSH session: %s\n", err)
		return nil, err
	}
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.ECHOCTL:       0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	if err := session.RequestPty(termType, height, width, modes); err != nil {
		fmt.Printf("request for pseudo terminal failed: %s\n", err)
		return nil, err
	}

	serverIn, err := session.StdinPipe()
	if err != nil {
		return nil, err
	}
	serverOut, err := session.StdoutPipe()
	if err != nil {
		return nil, err
	}
	eb, err := session.StderrPipe()
	if err != nil {
		return nil, err
	}
	clientIn, stdinPipe := io.Pipe() // You can treat stdinPipe as session.StdinPipe()
	stdoutPipe, clientOut := io.Pipe()
	// 设置trzsz
	trzszFilter := trzsz.NewTrzszFilter(clientIn, clientOut, serverIn, serverOut,
		trzsz.TrzszOptions{TerminalColumns: int32(width)})
	// 启动一个 shell 会话
	err = session.Shell()
	if err != nil {
		fmt.Printf("Failed to start SSH shell: %s\n", err)
		return nil, err
	}
	return &sshSession{
		client:      sshClient,
		session:     session,
		stdin:       stdinPipe,
		stdout:      stdoutPipe,
		stderr:      eb,
		trzszFilter: trzszFilter,
		clear: func() {
			defer func() {
				if r := recover(); r != nil {
					return
				}
			}()
			_ = sshClient.Close()
			_ = session.Close()
			_ = clientIn.Close()
			_ = stdinPipe.Close()
			_ = stdoutPipe.Close()
			_ = clientOut.Close()
		},
	}, nil
}
