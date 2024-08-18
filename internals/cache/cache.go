package cache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	Data map[string]cacheEntry
	sync.RWMutex
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		Data: make(map[string]cacheEntry),
	}

	cache.reapLoop(interval)

	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.Lock()
	defer c.Unlock()

	c.Data[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.RLock()
	defer c.RUnlock()

	ce, ok := c.Data[key]
	if !ok {
		return nil, false
	}

	return ce.val, true

}

func (c *Cache) reapLoop(interval time.Duration) {
	tick := time.Tick(interval)
	go func() {
		for range tick {
			base := time.Now()

			c.Lock()
			for k, v := range c.Data {
				if base.After(v.createdAt.Add(interval)) {
					delete(c.Data, k)
				}
			}
			c.Unlock()
		}
	}()

}
