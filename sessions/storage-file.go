package sessions

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"syscall"
)

//文件session仓库
type FileStorage struct {
	storagePath string
	prefix      string
	list        map[string]string
	rwLock      sync.RWMutex
}

func (fs *FileStorage) Save(w http.ResponseWriter, r *http.Request, sess *Session) error {
	fs.rwLock.Lock()
	defer fs.rwLock.Unlock()

	name := sess.ID
	filename := fs.prefix + name
	data, err := json.Marshal(sess)
	if err != nil {
		return err
	}
	err = fs.writeSessionFile(filename, string(data))
	if err != nil {
		return err
	}
	if sess.IsNew {
		sess.IsNew = false
		fs.list[name] = filename
		http.SetCookie(w, NewCookie(sess))
	}
	return nil
}

func (fs *FileStorage) Get(r *http.Request, name string) (*Session, error) {
	fs.rwLock.RLock()
	defer fs.rwLock.RUnlock()
	if sess_name, ok := fs.list[name]; ok {
		content, err := fs.readSessionFile(sess_name)
		if err != nil {
			return nil, err
		}
		session := &Session{}
		err = json.Unmarshal([]byte(content), &session)
		if err != nil {
			return nil, err
		}
		return session, nil
	}
	return nil, errors.New("session lost")
}

func (fs *FileStorage) Del(name string) {
	fs.rwLock.Lock()
	defer fs.rwLock.Unlock()
	delete(fs.list, name)
}
func (fs *FileStorage) GC() {

}

func (fs *FileStorage) readSessionFile(name string) (string, error) {
	file, err := os.Open(fs.storagePath + "/" + name)
	defer file.Close()
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
func (fs *FileStorage) writeSessionFile(name, content string) error {
	filename := fs.storagePath + "/" + name

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		_, err := os.Create(filename)

		if err != nil {
			return err
		}
	}

	file, err := os.OpenFile(filename, os.O_TRUNC, 0666)
	file.Close()
	_, err = io.WriteString(file, content)
	return err
}

func NewFileSessionStorage(path string, prefix ...string) {
	var sessionPrefix string
	var err error

	if err = syscall.Access(path, syscall.O_RDWR); err != nil {
		panic(err.Error())
	}

	if len(prefix) > 0 {
		sessionPrefix = prefix[0]
	} else {
		sessionPrefix = "sess_"
	}
	//判断路径是否可写
	file, err := os.Stat(path)
	if err != nil {
		panic(err.Error())
	}
	if !file.IsDir() {
		panic("session store path is not directory")
	}

	storage = &FileStorage{storagePath: path, prefix: sessionPrefix, list: make(map[string]string, 100)}
	storage.GC()
}
