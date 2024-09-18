package utils

import (
	"fmt"
	"github.com/bytedance/sonic"
	"os"
	"sync"
	"time"
)

type Env struct {
	Shell string `json:"SHELL"`
	Term  string `json:"TERM"`
}

type Header struct {
	Title     string `json:"title"`
	Version   int    `json:"version"`
	Height    int    `json:"height"`
	Width     int    `json:"width"`
	Env       Env    `json:"env"`
	Timestamp int    `json:"Timestamp"`
}

type Recorder struct {
	f   *os.File
	t   time.Time
	mux sync.Mutex
}

func (r *Recorder) Close() {
	if r.f != nil {
		_ = r.f.Close()
	}
}

func (r *Recorder) setHeader(header *Header) (err error) {
	var p []byte

	if p, err = sonic.Marshal(header); err != nil {
		return
	}
	r.mux.Lock()
	defer r.mux.Unlock()
	if _, err := r.f.Write(p); err != nil {
		return err
	}
	if _, err := r.f.Write([]byte("\n")); err != nil {
		return err
	}

	r.t = time.Unix(int64(header.Timestamp), 0)

	return
}

func (r *Recorder) Resize(rows, cols int) (n int, err error) {
	delta := time.Since(r.t).Seconds()
	jsonStr := fmt.Sprintf(`[%.6f, "r", "%dx%d"]\n`, delta, cols, rows)
	r.mux.Lock()
	defer r.mux.Unlock()
	if n, err = r.f.Write([]byte(jsonStr)); err != nil {
		return
	}
	return
}

func (r *Recorder) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	delta := time.Since(r.t).Seconds()
	var s []byte
	if s, err = sonic.Marshal(string(p)); err != nil {
		return
	}
	jsonStr := fmt.Sprintf(`[%.6f, "o", %s]\n`, delta, s)
	r.mux.Lock()
	defer r.mux.Unlock()
	if n, err = r.f.Write([]byte(jsonStr)); err != nil {
		return
	}
	return
}

func NewRecorder(recordingPath string, rows, cols int) (*Recorder, error) {
	if FileExists(recordingPath) {
		return nil, os.ErrExist
	}
	file, err := os.Create(recordingPath)
	if err != nil {
		return nil, err
	}

	recorder := &Recorder{}

	recorder.f = file

	recorder.t = time.Now()

	header := &Header{
		Title:     "",
		Version:   2,
		Height:    rows,
		Width:     cols,
		Env:       Env{Shell: "/bin/bash", Term: "xterm-256color"},
		Timestamp: int(time.Now().Unix()),
	}

	if err := recorder.setHeader(header); err != nil {
		return nil, err
	}

	return recorder, nil
}
