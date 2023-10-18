package cachestore

import (
	"fmt"
	"testing"
)

func TestCacheStore(t *testing.T) {
	cs := NewCacheStore[string, string](10, 1)
	t.Log(cs.Cap())
	for i := 0; i < 15; i++ {
		t.Log(cs.Save(fmt.Sprintf("hello%d", i), fmt.Sprintf("world%d", i)))
	}
	t.Log(cs.Len())
	t.Log(cs.Get("hello0"))
	t.Log(cs.Get("hello14"))

	cs.Save("mapvalue", "hello map")
	t.Log(cs.Get("mapvalue"))
}

func BenchmarkCasheStore(b *testing.B) {
	cs := NewCacheStore[string, interface{}](1000, 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cs.Save(fmt.Sprintf("hello%d", i), fmt.Sprintf("world%d", i))
		cs.Get(fmt.Sprintf("hello%d", i))
	}
}

func TestCacheStoreStruct(t *testing.T) {

	type A struct {
		Val string
	}
	cs := NewCacheStore[string, A](10, 1)
	t.Log(cs.Cap())
	for i := 0; i < 15; i++ {
		cs.Save(fmt.Sprintf("hello%d", i), A{Val: fmt.Sprintf("world%d", i)})
	}
	t.Log(cs.Len())
	a := cs.Get("hello10")
	t.Logf("%p\n", &a)
	a = cs.Get("hello10")
	t.Logf("%p\n", &a)

}
