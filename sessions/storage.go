package sessions

import (
	"errors"
	"net/http"
	"time"
)

//Storage Session仓库接口
//Get Save Del GC
type Storage interface {
	Save(http.ResponseWriter, *http.Request, *Session) error
	Get(*http.Request, string) (*Session, error)
	Del(string)
	GC()
}

//session仓库
var storage Storage

//NewSession 生成一个新的session类
func NewSession(path, domain string, maxage int, secure, httponly bool) *Session {
	cookieOpt := &CookieOptions{path, domain, maxage, secure, httponly}
	id := createSessionID()
	return &Session{
		ID:      id,
		Options: cookieOpt,
		storage: storage,
		Values:  make(map[interface{}]interface{}),
		IsNew:   true,
		ActTime: time.Now().Unix(),
	}
}

//GetSession 从请求中获取session
func GetSession(r *http.Request) (*Session, error) {
	cookie, err := r.Cookie(cookieSessionName)
	if err == nil {
		return storage.Get(r, cookie.Value)
	}
	return nil, errors.New("no session")
}

//DelSession 主动删除session
func DelSession(w http.ResponseWriter, sess *Session) {
	storage.Del(sess.ID)
	sess.Options.MaxAge = -1
	http.SetCookie(w, NewCookie(sess))
}

//CunstomSessionStorage 自定义存储引擎
func CunstomSessionStorage(store Storage) {
	storage = store
}
