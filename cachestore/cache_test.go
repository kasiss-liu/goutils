package cachestore

import (
	"fmt"
	"testing"
)

func TestCacheStore(t *testing.T) {
	cs := NewCacheStore(10, 1)
	t.Log(cs.Cap())
	for i := 0; i < 15; i++ {
		t.Log(cs.Save(fmt.Sprintf("hello%d", i), fmt.Sprintf("world%d", i)))
	}
	t.Log(cs.Len())
	t.Log(cs.Get("hello0"))
	t.Log(cs.Get("hello14"))

	cs.Save("mapvalue", map[string]string{"1": "aa"})
	t.Log(cs.Get("mapvalue"))
}

func BenchmarkCasheStore(b *testing.B) {
	cs := NewCacheStore(1000, 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cs.Save(fmt.Sprintf("hello%d", i), fmt.Sprintf("world%d", i))
		cs.Get(fmt.Sprintf("hello%d", i))
	}
}
