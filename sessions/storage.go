package sessions

import (
	"errors"
	"net/http"
	"sync"
	"time"
)

//Session仓库接口
//Get Save Del GC
type Storage interface {
	Save(http.ResponseWriter, *http.Request, *Session) error
	Get(*http.Request, string) (*Session, error)
	Del(string)
	GC()
}

var memStorage Storage

func init() {
	memStorage = MemStorage{list: make(map[string]*Session, 100)}
	memStorage.GC()
}

//实现一个仓库接口
//session存储于内存中，服务重启后丢失
type MemStorage struct {
	list   map[string]*Session
	rwLock sync.RWMutex
}

//仓库save方法 将sessionId数据写入cookie返回到客户端
//并将session内容存入内存仓库
func (m MemStorage) Save(w http.ResponseWriter, r *http.Request, sess *Session) error {
	m.rwLock.Lock()
	defer m.rwLock.Unlock()

	name := sess.ID
	m.list[name] = sess
	if sess.IsNew {
		sess.IsNew = false
		http.SetCookie(w, NewCookie(sess))
	}
	return nil

}

//从仓库中获取一个session
func (m MemStorage) Get(r *http.Request, name string) (*Session, error) {
	m.rwLock.RLock()
	defer m.rwLock.RUnlock()
	if sess, ok := m.list[name]; ok {
		return sess, nil
	}
	return nil, errors.New("session lost")
}

//从仓库中删除一个session
func (m MemStorage) Del(name string) {
	m.rwLock.Lock()
	defer m.rwLock.Unlock()
	delete(m.list, name)
}

//仓库内过期session清除
//每秒钟筛选一遍
func (m MemStorage) GC() {
	go func() {
		for {
			for _, session := range m.list {
				if session.GC() {
					m.rwLock.Lock()
					m.Del(session.ID)
					m.rwLock.Unlock()
				}
			}
			time.Sleep(1 * time.Second)
		}
	}()

}

//生成一个新的session类
func NewSession(path, domain string, maxage int, secure, httponly bool) *Session {
	cookieOpt := &CookieOptions{path, domain, maxage, secure, httponly}
	id := createSessionId()
	return &Session{
		ID:      id,
		Options: cookieOpt,
		storage: memStorage,
		Values:  make(map[interface{}]interface{}),
		IsNew:   true,
		ActTime: time.Now().Unix(),
	}
}

//从请求中获取session
func GetSession(r *http.Request) (*Session, error) {
	cookie, err := r.Cookie(CookieSessionName)
	if err == nil {
		return memStorage.Get(r, cookie.Value)
	}
	return nil, errors.New("no session")
}

//主动删除session
func DelSession(w http.ResponseWriter, sess *Session) {
	memStorage.Del(sess.ID)
	sess.Options.MaxAge = -1
	http.SetCookie(w, NewCookie(sess))
}
