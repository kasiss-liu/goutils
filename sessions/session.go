package sessions

import (
	"bytes"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

//cookie内存放的sessionId 键名
var (
	cookieSessionName string = "GO_WEBSESS"
)

//CookieOptions
//cookie存放的基础属性
//路径、所属域、存活时间、是否安全、只经由http传输
type CookieOptions struct {
	Path     string
	Domain   string
	MaxAge   int
	Secure   bool
	HttpOnly bool
}

//Session 结构
//id、值、cookie属性、是否是新会话、最后活跃时间、仓库
type Session struct {
	ID      string
	Values  map[interface{}]interface{}
	Options *CookieOptions
	storage Storage
	IsNew   bool
	ActTime int64
}

//设置session值
func (s *Session) Set(key interface{}, value interface{}) {
	s.Values[key] = value
}

//获取session内的值
//获取后需要自行断言
func (s *Session) Get(key interface{}) interface{} {
	var value interface{}
	if value, ok := s.Values[key]; ok {
		return value
	}
	return value
}

//删除某个session值
func (s *Session) Del(key interface{}) {
	if _, ok := s.Values[key]; ok {
		delete(s.Values, key)
	}
}

//获取一个session中值的个数
func (s *Session) Len() (n int) {
	n = len(s.Values)
	return
}

//将session保存
func (s *Session) Save(w http.ResponseWriter, r *http.Request) {
	s.ActTime = time.Now().Unix()
	s.storage.Save(w, r, s)
}

//session 垃圾回收判断
func (s *Session) GC() bool {
	return int(time.Now().Unix()-s.ActTime) > s.Options.MaxAge
}

//生成一个新的Cookie结构
func NewCookie(s *Session) *http.Cookie {
	cookie := &http.Cookie{
		Name:   cookieSessionName,
		Value:  s.ID,
		Path:   s.Options.Path,
		Domain: s.Options.Domain,
		Secure: s.Options.Secure,
		MaxAge: s.Options.MaxAge,
	}
	return cookie
}

//生成随机sessionId
func createSessionId() string {

	var result bytes.Buffer
	var temp string
	for i := 0; i < 10; {
		temp = getChar()
		result.WriteString(temp)
		i++

	}
	return result.String()
}

//获取随机字符串
func getChar() string {
	rand.Seed(time.Now().UnixNano())
	switch rand.Intn(3) {
	case 0:
		return string(65 + rand.Intn(90-65))
	case 1:
		return string(97 + rand.Intn(122-97))
	default:
		return strconv.Itoa(rand.Intn(9))
	}
}

func Init(store Storage, cookieName ...string) {
	if len(cookieName) > 0 {
		cookieSessionName = cookieName[0]
	}
	storage = store
	storage.GC()
}

//设置cookieSessionName 来取代默认值
func SetCookieSessionName(s string) {
	cookieSessionName = s
}
