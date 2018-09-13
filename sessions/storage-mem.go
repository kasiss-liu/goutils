package sessions

import (
	"errors"
	"net/http"
	"sync"
	"time"
)

//MemStorage 实现一个仓库接口
//session存储于内存中，服务重启后丢失
type MemStorage struct {
	list   map[string]*Session
	rwLock sync.RWMutex
}

//Save 仓库save方法 将sessionId数据写入cookie返回到客户端
//并将session内容存入内存仓库
func (ms *MemStorage) Save(w http.ResponseWriter, r *http.Request, sess *Session) error {
	ms.rwLock.Lock()
	defer ms.rwLock.Unlock()

	name := sess.ID
	ms.list[name] = sess
	if sess.IsNew {
		sess.IsNew = false
		http.SetCookie(w, NewCookie(sess))
	}
	return nil

}

//Get 从仓库中获取一个session
func (ms *MemStorage) Get(r *http.Request, name string) (*Session, error) {
	ms.rwLock.RLock()
	defer ms.rwLock.RUnlock()
	if sess, ok := ms.list[name]; ok {
		return &Session{
			ID:      sess.ID,
			Values:  sess.Values,
			Options: sess.Options,
			IsNew:   false,
			ActTime: sess.ActTime,
			storage: storage,
		}, nil
	}
	return nil, errors.New("session lost")
}

//Del 从仓库中删除一个session
func (ms *MemStorage) Del(name string) {
	ms.rwLock.Lock()
	defer ms.rwLock.Unlock()
	delete(ms.list, name)
}

//GC 仓库内过期session清除
//每秒钟筛选一遍
func (ms *MemStorage) GC() {
	go func() {
		for {
			for _, session := range ms.list {
				if session.GC() {
					ms.rwLock.Lock()
					ms.Del(session.ID)
					ms.rwLock.Unlock()
				}
			}
			time.Sleep(1 * time.Second)
		}
	}()

}

//NewMemSessionStorage 生成一个新的内存session仓库
func NewMemSessionStorage() Storage {
	return &MemStorage{list: make(map[string]*Session, 100)}
}
