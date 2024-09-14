package utils

import (
	"github.com/bytedance/sonic"
	"os"
	"sync/atomic"
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
	File    *os.File
	Time    time.Time
	Enabled *atomic.Bool
}

func (recorder *Recorder) Enable() {
	recorder.Enabled.Store(true)
}
func (recorder *Recorder) Disable() {
	recorder.Enabled.Store(false)
	recorder.Close()
}

func (recorder *Recorder) Close() {
	recorder.Enabled.Store(false)
	if recorder.File != nil {
		_ = recorder.File.Close()
	}
}

func (recorder *Recorder) setHeader(header *Header) (err error) {
	var p []byte

	if p, err = sonic.Marshal(header); err != nil {
		return
	}

	if _, err := recorder.File.Write(p); err != nil {
		return err
	}
	if _, err := recorder.File.Write([]byte("\n")); err != nil {
		return err
	}

	recorder.Time = time.Unix(int64(header.Timestamp), 0)

	return
}

func (recorder *Recorder) Write(p []byte) (n int, err error) {
	delta := time.Since(recorder.Time).Seconds()

	row := make([]interface{}, 0)
	row = append(row, delta)
	row = append(row, "o")
	row = append(row, string(p))

	var s []byte
	if s, err = sonic.Marshal(row); err != nil {
		return
	}
	if n, err = recorder.File.Write(s); err != nil {
		return
	}
	_, err = recorder.File.Write([]byte("\n"))
	return
}

func NewRecorder(recordingPath string) (recorder *Recorder, err error) {
	recorder = &Recorder{
		Enabled: &atomic.Bool{},
	}
	parentDirectory := GetParentDirectory(recordingPath)
	if FileExists(parentDirectory) {
		if err := os.RemoveAll(parentDirectory); err != nil {
			return nil, err
		}
	}

	if err = os.MkdirAll(parentDirectory, 0777); err != nil {
		return
	}

	var file *os.File
	file, err = os.Create(recordingPath)
	if err != nil {
		return nil, err
	}

	recorder.Enabled.Store(true)
	recorder.File = file

	header := &Header{
		Title:     "",
		Version:   2,
		Height:    40,
		Width:     80,
		Env:       Env{Shell: "/bin/bash", Term: "xterm-256color"},
		Timestamp: int(time.Now().Unix()),
	}

	if err := recorder.setHeader(header); err != nil {
		return nil, err
	}

	return recorder, nil
}
