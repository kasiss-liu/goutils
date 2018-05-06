package fluentLogger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"
)

type IWriter interface {
	Write([]byte) (int, error)
	WriteString(string) (int, error)
}

//三种日志命名方式
//不分割
//按天分割
//按小时分割
const (
	PARTITION_NONE = iota
	PARTITION_DAY
	PARTITION_HOUR
)

//日志写入器
type FileWriter struct {
	fLock  sync.Mutex
	prefix string
	path   string
	mode   int
	file   *os.File
	fname  string
}

func (fw *FileWriter) Write(bs []byte) (int, error) {
	fw.fLock.Lock()
	defer fw.fLock.Unlock()
	return io.WriteString(fw.file, string(bs))
}

func (fw *FileWriter) WriteString(s string) (int, error) {
	fw.fLock.Lock()
	defer fw.fLock.Unlock()
	return io.WriteString(fw.file, s)
}

//每次写入前校验一下是否被分割
func (fw *FileWriter) Prepare() error {
	//如果文件地址未改变 则不处理
	if fw.ensureFileName() {
		return nil
	}
	//如果文件资源存在 则进行关闭
	if fw.file != nil {
		fw.file.Close()
	}
	//打开新的文件
	f, err := os.OpenFile(fw.fname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	fw.file = f
	return nil
}

//每次校对是否该进行文件分割
func (fw *FileWriter) ensureFileName() bool {
	//计算新的文件名
	var filename string
	switch fw.mode {
	case PARTITION_NONE:
		filename = fw.path + "/" + fw.prefix + ".log"
	case PARTITION_DAY:
		y, m, d := time.Now().Date()
		filename = fw.path + "/" + fw.prefix + fmt.Sprintf("_%d%02d%02d", y, m, d) + ".log"
	case PARTITION_HOUR:
		t := time.Now()
		y, m, d := t.Date()
		h := t.Hour()
		filename = fw.path + "/" + fw.prefix + fmt.Sprintf("_%d%02d%02d%02d", y, m, d, h) + ".log"
	default:
		filename = fw.path + "/" + fw.prefix + ".log"
	}
	//如果发生改变则更改日志文件名称
	same := filename == fw.fname
	if !same {
		fw.fname = filename
	}
	return same
}

//主动关闭文件资源
func (fw *FileWriter) Close() error {
	return fw.file.Close()
}

func NewFileWriter(dir, prefix string, mode int) *FileWriter {
	err := syscall.Access(dir, syscall.O_RDWR)
	if err != nil {
		panic(dir + " have no access\nErrorInfo:" + err.Error())
	}
	file := &FileWriter{
		prefix: prefix,
		mode:   mode,
		path:   strings.TrimRight(dir, "/"),
	}
	file.Prepare()
	return file
}
