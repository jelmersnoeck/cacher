// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package cacher

import (
	"runtime"
	"time"
	"unsafe"
)

type cachedItem struct {
	value  interface{}
	expiry time.Time
	expire bool
}

// MemoryCache is a caching implementation that stores the data in memory. The
// cache will be emptied when the application has run.
type MemoryCache struct {
	items map[string]cachedItem
	keys  []string
	limit uintptr
	size  uintptr
}

// NewMemoryCache creates a new instance of MemoryCache and initiates the
// storage map.
func NewMemoryCache(limit uintptr) *MemoryCache {
	cache := new(MemoryCache)
	cache.items = make(map[string]cachedItem)
	if limit == 0 {
		// 10% of system memory
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		cache.limit = uintptr(float64(memStats.Sys) * 0.1)
	} else {
		cache.limit = limit
	}
	cache.size = 0

	return cache
}

// Add an item to the cache. If the item is already cached, the value won't be
// overwritten.
//
// ttl defines the number of seconds the value should be cached. If ttl is 0,
// the item will be cached infinitely.
func (c *MemoryCache) Add(key string, value interface{}, ttl int) bool {
	if c.exists(key) {
		return false
	}

	return c.Set(key, value, ttl)
}

// Set sets the value of an item, regardless of wether or not the value is
// already cached.
//
// ttl defines the number of seconds the value should be cached. If ttl is 0,
// the item will be cached infinitely.
func (c *MemoryCache) Set(key string, value interface{}, ttl int) bool {
	expiry := time.Now().Add(time.Duration(ttl) * time.Second)

	var expire bool
	if ttl > 0 {
		expire = true
	}

	c.items[key] = cachedItem{value, expiry, expire}
	c.keys = append(c.keys, key)
	c.size += unsafe.Sizeof(c.items[key])
	c.lru(key)
	c.evict()
	return true
}

// Get gets the value out of the map associated with the provided key.
func (c *MemoryCache) Get(key string) (interface{}, bool) {
	if c.exists(key) {
		return c.items[key].value, true
	}

	return nil, false
}

// Flush will remove all the items from the hash.
func (c *MemoryCache) Flush() bool {
	c.items = make(map[string]cachedItem)
	c.size = 0
	return true
}

// Delete will validate if the key actually is stored in the cache. If it is
// stored, it will remove the item from the cache. If it is not stored, it will
// return false.
func (c *MemoryCache) Delete(key string) bool {
	for i, v := range c.keys {
		if v == key {
			return c.removeAt(i)
		}
	}

	return false
}

func (c *MemoryCache) removeAt(index int) bool {
	key := c.keys[index]
	c.keys = append(c.keys[:index], c.keys[index+1:]...)
	c.size -= unsafe.Sizeof(c.items[key])
	delete(c.items, key)

	return true
}

// exists checks if a key is stored in the cache.
//
// If the key is stored in the cache, but the expiry date has passed, we will
// remove the item from the cache and return false. If the expiry has not passed
// yet, it will return false.
func (c *MemoryCache) exists(key string) bool {
	cachedItem, exists := c.items[key]
	if exists {
		if !cachedItem.expire || time.Now().Before(cachedItem.expiry) {
			c.lru(key)
			return true
		}

		// Item is expired, delete it and act as it doesn't exist
		c.Delete(key)
	}

	return false
}

// evict clears off the items in the cache that have been least active.
func (c *MemoryCache) evict() {
	for {
		if c.size > c.limit {
			c.removeAt(0)
		} else {
			break
		}
	}
}

// lru stands for Least Recently Used. We will use this algorithm to mark items
// that are not active in our cache to be freed when the size is over its limit.
func (c *MemoryCache) lru(key string) {
	for i, v := range c.keys {
		if v == key {
			c.keys = append(c.keys[:i], c.keys[i+1:]...)
			c.keys = append(c.keys, key)
		}
	}
}
