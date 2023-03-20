package geecache

import (
	"sanHeRecruitment/module/lruEngine/lru"
	"sync"
)

type LruCache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

func (c *LruCache) Add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *LruCache) Get(key string) (value ByteView, ok bool) {
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

func (c *LruCache) Delete(key string) (ok bool) {
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
