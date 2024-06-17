package termx

import (
	"github.com/pkg/sftp"
)

type PtyX interface {
	Resize(rows, cols int) error
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	Close() error
	Sftp() (*sftp.Client, error)
}
