package cachestore

import (
	"sync"
	"time"
)

//CacheStore 缓存仓
type CacheStore[K comparable, V any] struct {
	data     map[K]*cacheItem[K, V]
	gcs      []*cacheItem[K, V]
	cap      int
	lifetime int
	lock     sync.Mutex
}

type cacheItem[K comparable, V any] struct {
	key    K
	value  V
	gcTime time.Time
}

//Len 获取缓存仓当前存储量
func (cs *CacheStore[K, V]) Len() int {
	return len(cs.data)
}

//Cap 获取缓存仓容量
func (cs *CacheStore[K, V]) Cap() int {
	return cs.cap
}

//Save 保存一个元素值
func (cs *CacheStore[K, V]) Save(key K, value V) bool {

	cs.lock.Lock()
	defer cs.lock.Unlock()
	cs.toCheckOrGcFirstItem()
	item := &cacheItem[K, V]{
		key:    key,
		value:  value,
		gcTime: time.Now().Add(time.Second * time.Duration(cs.lifetime)),
	}
	cs.data[key] = item
	cs.gcs = append(cs.gcs, item)
	return true
}

//Get  获取一个元素值
func (cs *CacheStore[K, V]) Get(key K) V {
	var res V
	cs.lock.Lock()
	defer cs.lock.Unlock()
	if v, ok := cs.data[key]; ok {
		if time.Now().Before(v.gcTime) {
			return v.value
		}
	}
	return res
}

//容量检测，去除最早的元素
func (cs *CacheStore[K, V]) toCheckOrGcFirstItem() {
	if len(cs.data) >= cs.cap {
		item := cs.gcs[0]
		if item != nil {
			delete(cs.data, item.key)
			cs.gcs = cs.gcs[1:]
		}
	}
}

//过期回收
func (cs *CacheStore[K, V]) gc() {
	for {
		<-time.Tick(1 * time.Second)
		cs.doGC()
	}
}

//单独进行一次回收处理
func (cs *CacheStore[K, V]) doGC() {
	grayKey := make(map[int]bool, 10)
	for k, item := range cs.gcs {

		if item != nil && time.Now().After(item.gcTime) {
			cs.lock.Lock()
			delete(cs.data, item.key)
			cs.lock.Unlock()
			grayKey[k] = true
		}
	}
	cs.lock.Lock()
	newGcs := make([]*cacheItem[K, V], 0, cap(cs.gcs))
	for k, item := range cs.gcs {
		if _, ok := grayKey[k]; !ok {
			newGcs = append(newGcs, item)
		}
	}
	cs.gcs = newGcs
	cs.lock.Unlock()
}

//NewCacheStore 获取一个新的缓存仓
func NewCacheStore[K comparable, V any](cap, lifetime int) *CacheStore[K, V] {
	cs := &CacheStore[K, V]{
		data:     make(map[K]*cacheItem[K, V], cap),
		gcs:      make([]*cacheItem[K, V], 0, cap),
		cap:      cap,
		lifetime: lifetime,
	}
	go cs.gc()
	return cs
}
