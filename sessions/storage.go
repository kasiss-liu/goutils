package sessions

import (
	"errors"
	"net/http"
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

//session仓库
var storage Storage

//生成一个新的session类
func NewSession(path, domain string, maxage int, secure, httponly bool) *Session {
	cookieOpt := &CookieOptions{path, domain, maxage, secure, httponly}
	id := createSessionId()
	return &Session{
		ID:      id,
		Options: cookieOpt,
		storage: storage,
		Values:  make(map[interface{}]interface{}),
		IsNew:   true,
		ActTime: time.Now().Unix(),
	}
}

//从请求中获取session
func GetSession(r *http.Request) (*Session, error) {
	cookie, err := r.Cookie(CookieSessionName)
	if err == nil {
		return storage.Get(r, cookie.Value)
	}
	return nil, errors.New("no session")
}

//主动删除session
func DelSession(w http.ResponseWriter, sess *Session) {
	storage.Del(sess.ID)
	sess.Options.MaxAge = -1
	http.SetCookie(w, NewCookie(sess))
}

func CunstomSessionStorage(store Storage) {
	storage = store
}
