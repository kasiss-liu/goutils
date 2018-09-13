package fluentLogger

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

//日志结构
type logData struct {
	dataType string
	content  string
	time     time.Time
}

//定义一些必要的参数
var (
	loggers           map[string]Logger //写入器池
	DefaultLogChanLen = 100             //默认日志数据通道长度
	DefaultLoggers    = 5               //默认写入器池初始大小
	logLock           sync.WaitGroup    //同步锁
)

//初始化写入器池
func init() {
	loggers = make(map[string]Logger, DefaultLoggers)
}

//Logger 写入器结构
type Logger struct {
	writer     IWriter
	logChannel chan logData
}

//Log 普通log数据记录 向通道内写入数据
func (l *Logger) Log(dataType, content string, t time.Time) error {
	l.logChannel <- logData{dataType: dataType, content: content, time: t}
	//写入一次 计数器+1
	logLock.Add(1)
	return nil
}

//数据通道监控goroutiue
func (l *Logger) monitorChannel() {
	for {
		c := <-l.logChannel
		go l.writeLog(c)
	}
}

//数据结构 数据写入
func (l *Logger) writeLog(c logData) {
	switch l.writer.(type) {
	case *FileWriter:
		content := fmt.Sprintf("[%s]%s:%s\r\n", c.time.Format("2006-01-02 15:04:05"), c.dataType, c.content)
		l.writer.(*FileWriter).Prepare()
		l.writer.WriteString(content)
	}
	//执行一次计数器-1
	logLock.Done()
}

//Log 外部调用日志
func Log(mark, content string, t time.Time) error {
	if logger, ok := loggers[mark]; ok {
		return logger.Log("Log", content, t)
	}
	return errors.New("Wrong Mark: " + mark)
}

//Notice 外部调用notice
func Notice(mark, content string, t time.Time) error {
	if logger, ok := loggers[mark]; ok {
		return logger.Log("Notice", content, t)
	}
	return errors.New("Wrong Mark: " + mark)
}

//Warning 外部调用warning
func Warning(mark, content string, t time.Time) error {
	if logger, ok := loggers[mark]; ok {
		return logger.Log("Warning", content, t)
	}
	return errors.New("Wrong Mark: " + mark)
}

//Error 外部调用error
func Error(mark, content string, t time.Time) error {
	if logger, ok := loggers[mark]; ok {
		return logger.Log("Error", content, t)
	}
	return errors.New("Wrong Mark: " + mark)
}

//Write 自定义写入前缀
func Write(mark, prefix, content string, t time.Time) error {
	if logger, ok := loggers[mark]; ok {
		return logger.Log(prefix, content, t)
	}
	return errors.New("Wrong Mark: " + mark)
}

//RegisterFileLogger 注册一个新的文件写入器到池
func RegisterFileLogger(mark, dir, prefix string, mode int) bool {
	//验证写入器是否已经存在
	if _, ok := loggers[mark]; ok {
		return false
	}
	loggers[mark] = Logger{
		writer:     NewFileWriter(dir, prefix, mode),
		logChannel: make(chan logData, DefaultLogChanLen),
	}
	return true
}

//StartLogger 启动日志写入器
func StartLogger() error {
	if len(loggers) < 1 {
		return errors.New("No loggers registered ! ")
	}
	for _, v := range loggers {
		go v.monitorChannel()
	}
	return nil
}
