package lruEngine

import (
	"sanHeRecruitment/library/lruEngine/lru"
	"sync"
)

type lruCache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

var LruEngine = lruCache{
	lru: lru.New(2<<20*10, nil),
}

func (c *lruCache) Add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}
	c.lru.Add(key, value)
}

func (c *lruCache) Get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}

	if val, ok := c.lru.Get(key); ok {
		return val.(ByteView), ok
	}
	return
}

func (c *lruCache) Delete(key string) (ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}

	if ok := c.lru.PreRemove(key); ok {
		return ok
	}
	return false
}
