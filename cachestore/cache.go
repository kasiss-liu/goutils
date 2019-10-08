package cachestore

import (
	"sync"
	"time"
)

type cacheKey interface{}

//CacheStore 缓存仓
type CacheStore struct {
	data     map[cacheKey]*cacheItem
	gcs      []*cacheItem
	cap      int
	lifetime int
	lock     sync.Mutex
}

//cacheItem 缓存仓元素
type cacheItem struct {
	key    cacheKey
	value  interface{}
	gcTime time.Time
}

//Len 获取缓存仓当前存储量
func (cs *CacheStore) Len() int {
	return len(cs.data)
}

//Cap 获取缓存仓容量
func (cs *CacheStore) Cap() int {
	return cs.cap
}

//Save 保存一个元素值
func (cs *CacheStore) Save(key, value interface{}) bool {

	k := cacheKey(key)
	cs.lock.Lock()
	defer cs.lock.Unlock()
	cs.toCheckOrGcFirstItem()
	item := &cacheItem{
		key:    k,
		value:  value,
		gcTime: time.Now().Add(time.Second * time.Duration(cs.lifetime)),
	}
	cs.data[k] = item
	cs.gcs = append(cs.gcs, item)
	return true
}

//Get  获取一个元素值
func (cs *CacheStore) Get(key interface{}) interface{} {
	k := cacheKey(key)
	cs.lock.Lock()
	defer cs.lock.Unlock()
	if v, ok := cs.data[k]; ok {
		if time.Now().Before(v.gcTime) {
			return v.value
		}
	}
	return ""
}

//容量检测，去除最早的元素
func (cs *CacheStore) toCheckOrGcFirstItem() {
	if len(cs.data) >= cs.cap {
		item := cs.gcs[0]
		if item != nil {
			delete(cs.data, item.key)
			cs.gcs = cs.gcs[1:]
		}
	}
}

//过期回收
func (cs *CacheStore) gc() {
	for {
		<-time.Tick(1 * time.Second)
		cs.doGC()
	}
}
//单独进行一次回收处理
func (cs *CacheStore) doGC() {
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
	newGcs := make([]*cacheItem,0,cap(cs.gcs))
	for k, item := range cs.gcs {
		if _, ok := grayKey[k]; !ok {
			newGcs = append(newGcs, item)
		}
	}
	cs.gcs = newGcs
	cs.lock.Unlock()
}

//NewCacheStore 获取一个新的缓存仓
func NewCacheStore(cap, lifetime int) *CacheStore {
	cs := &CacheStore{
		data:     make(map[cacheKey]*cacheItem, cap),
		gcs:      make([]*cacheItem, 0, cap),
		cap:      cap,
		lifetime: lifetime,
	}
	go cs.gc()
	return cs
}
