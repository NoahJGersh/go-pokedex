package pokeutils

import (
	"sync"
	"time"
)

type Cache struct {
	entries map[string]cacheEntry
	mutex   *sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) Cache {
	var mutex sync.Mutex
	c := Cache{
		entries: make(map[string]cacheEntry),
		mutex:   &mutex,
	}

	go c.reapLoop(interval)

	return c
}

func (c Cache) Add(key string, val []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c Cache) Get(key string) ([]byte, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, ok := c.entries[key]
	if !ok {
		return []byte{}, false
	}

	return entry.val, true
}

func (c Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		tickTime := <-ticker.C
		for k, v := range c.entries {
			if tickTime.Sub(v.createdAt) > interval {
				delete(c.entries, k)
			}
		}
	}
}
