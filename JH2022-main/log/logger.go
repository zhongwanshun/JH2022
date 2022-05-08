package log

import (
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"time"
)

var (
	Trace   *log.Logger //  任何信息都可以输出
	Info    *log.Logger //  重要信息
	Warning *log.Logger //  警告
	Error   *log.Logger //  错误
)

type FileLog struct {
	logPath string
	logName string
}

func NewLogger(p string) {
	logfile := &FileLog{
		logPath: path.Join(p, "logger", time.Now().Format("2006-01-02")),
		logName: strconv.Itoa(time.Now().Hour()) + ".log",
	}
	// 创建日志目录
	err := os.MkdirAll(logfile.logPath, 0777)
	if err != nil {
		log.Println("创建日志目录失败:", err)
		return
	}
	// 创建日志文件
	file, err := os.OpenFile(path.Join(logfile.logPath, logfile.logName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("创建日志文件失败:", err)
		return
	}
	// 创建日志输出对象
	Trace = log.New(io.MultiWriter(file, os.Stdout),
		"TRACE: ",
		log.Ldate|log.Lmicroseconds|log.Lshortfile)
	Info = log.New(io.MultiWriter(file, os.Stdout),
		"INFO: ",
		log.Ldate|log.Lmicroseconds|log.Lshortfile)
	Warning = log.New(io.MultiWriter(file, os.Stdout),
		"WARNING: ",
		log.Ldate|log.Lmicroseconds|log.Lshortfile)
	Error = log.New(io.MultiWriter(file, os.Stderr),
		"ERROR: ",
		log.Ldate|log.Lmicroseconds|log.Llongfile)
}
