package utils

import (
	"bufio"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

type LogFilter struct {
	logFile    *os.File
	inTransfer bool
	buffer     strings.Builder
}

func NewLogFilter(logFilePath string) (*LogFilter, error) {
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &LogFilter{logFile: file}, nil
}

func (lf *LogFilter) Close() {
	lf.logFile.Close()
}

func (lf *LogFilter) ProcessOutput(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		lf.processLine(line)
	}
	if err := scanner.Err(); err != nil {
		log.Printf("读取输出时出错: %v", err)
	}
}

func (lf *LogFilter) processLine(line string) {
	// 检测 trzsz 或 lrzsz 传输开始
	if strings.Contains(line, "rz waiting to receive.") ||
		strings.Contains(line, "Starting zmodem transfer") {
		lf.inTransfer = true
		lf.buffer.Reset()
		return
	}

	// 检测传输结束
	if lf.inTransfer && (strings.Contains(line, "Transfer complete") ||
		strings.Contains(line, "rz transfer complete")) {
		lf.inTransfer = false
		lf.writeLog("文件传输完成\n")
		return
	}

	// 如果不在传输中，写入日志
	if !lf.inTransfer {
		lf.writeLog(line + "\n")
	} else {
		// 在传输中，但可能包含进度信息
		if match, _ := regexp.MatchString(`\d+%`, line); match {
			lf.buffer.WriteString(line + "\n")
		}
	}
}

func (lf *LogFilter) writeLog(content string) {
	if _, err := lf.logFile.WriteString(content); err != nil {
		log.Printf("写入日志时出错: %v", err)
	}
}
