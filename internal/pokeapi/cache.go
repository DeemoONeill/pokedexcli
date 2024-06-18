package pokeapi

import (
	"sync"
	"time"
)

type Cache struct {
	cache    map[string]cacheEntry
	interval time.Duration
	lock     *sync.RWMutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(ttl time.Duration) Cache {
	cache := Cache{interval: ttl, cache: make(map[string]cacheEntry), lock: &sync.RWMutex{}}
	ticker := time.NewTicker(ttl)
	go cache.reapLoop(ticker)
	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.cache[key] = cacheEntry{time.Now(), val}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	value, ok := c.cache[key]
	return value.val, ok
}

func (c *Cache) reapLoop(tick *time.Ticker) {
	for range tick.C {
		for k, entry := range c.cache {
			if time.Since(entry.createdAt) > c.interval {
				delete(c.cache, k)
			}
		}
	}
}
