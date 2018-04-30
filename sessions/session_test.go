package sessions

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func initMem(store Storage) {

	Init(store, "TEST")
}

func TestMemSessions(t *testing.T) {
	store := NewMemSessionStorage()
	initMem(store)
	newSession := NewSession("/", "localhost", 300, true, false)

	newSession.Set("test", "test111")

	getValue := newSession.Get("test")
	fmt.Println("val:", getValue)

	l := newSession.Len()
	fmt.Println("len:", l)

	newSession.Del("test")

	l = newSession.Len()
	fmt.Println("len:", l)

	isGc := newSession.GC()
	fmt.Println("Gc:", isGc)

	id := newSession.ID
	name := cookieSessionName

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "localhost:8999", nil)

	newSession.Save(resp, nil)

	req.AddCookie(&http.Cookie{
		Name:     name,
		Value:    id,
		Path:     "/",
		Domain:   "localhost",
		HttpOnly: false,
		Secure:   false,
	})

	sess, err := GetSession(req)
	if err == nil {
		fmt.Println("getSess:", sess)
	} else {
		t.Error(err.Error())
	}

	DelSession(resp, newSession)

}

func TestFileSessions(t *testing.T) {
	store := NewFileSessionStorage("./", "sess_")
	initMem(store)
	CunstomSessionStorage(store)
	SetCookieSessionName("TEST_SESSION")
	newSession := NewSession("/", "localhost", 300, false, false)

	getValue := newSession.Get("test")
	fmt.Println("val:", getValue)

	l := newSession.Len()
	fmt.Println("len:", l)

	newSession.Del("test")

	l = newSession.Len()
	fmt.Println("len:", l)

	isGc := newSession.GC()
	fmt.Println("Gc:", isGc)

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "localhost:8999", nil)
	newSession.Save(resp, nil)

	id := newSession.ID
	name := cookieSessionName

	req.AddCookie(&http.Cookie{
		Name:     name,
		Value:    id,
		Path:     "/",
		Domain:   "localhost",
		HttpOnly: false,
		Secure:   false,
	})

	sess, err := GetSession(req)
	if err == nil {
		fmt.Println("getSess:", sess)
	} else {
		t.Error(err.Error())
	}

	DelSession(resp, newSession)
}
