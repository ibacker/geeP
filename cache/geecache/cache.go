package geecache

import (
	"geecache/lru"
	"sync"
)

type cache struct {
	mu  sync.Mutex
	lru *lru.Cache
	// 最大缓存大小
	cacheBytes int64
}

// add
func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

// get 返回缓存值
func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}

	if value, ok := c.lru.Get(key); ok {
		return value.(ByteView), ok
	}
	return
}
