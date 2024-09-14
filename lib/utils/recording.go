package utils

import (
	"github.com/bytedance/sonic"
	"os"
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
	File *os.File
	Time time.Time
}

func (recorder *Recorder) Close() {
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

func NewRecorder(recordingPath string) (*Recorder, error) {
	parentDirectory := GetParentDirectory(recordingPath)
	if FileExists(recordingPath) {
		if err := os.RemoveAll(recordingPath); err != nil {
			return nil, err
		}
	}

	if err := os.MkdirAll(parentDirectory, 0777); err != nil {
		return nil, err
	}

	var file *os.File
	file, err := os.Create(recordingPath)
	if err != nil {
		return nil, err
	}

	recorder := &Recorder{}

	recorder.File = file

	recorder.Time = time.Now()

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
