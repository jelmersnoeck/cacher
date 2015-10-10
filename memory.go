// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package cacher

import "time"

type cachedItem struct {
	value  interface{}
	expiry time.Time
	expire bool
}

// MemoryCache is a caching implementation that stores the data in memory. The
// cache will be emptied when the application has run.
type MemoryCache struct {
	items map[string]cachedItem
}

// NewMemoryCache creates a new instance of MemoryCache and initiates the
// storage map.
func NewMemoryCache() *MemoryCache {
	cache := new(MemoryCache)
	cache.items = make(map[string]cachedItem)

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
	return true
}

// Get gets the value out of the map associated with the provided key.
func (c *MemoryCache) Get(key string) interface{} {
	return c.items[key].value
}

// Flush will remove all the items from the hash.
func (c *MemoryCache) Flush() bool {
	c.items = make(map[string]cachedItem)
	return true
}

// Delete will validate if the key actually is stored in the cache. If it is
// stored, it will remove the item from the cache. If it is not stored, it will
// return false.
func (c *MemoryCache) Delete(key string) bool {
	_, exists := c.items[key]

	if exists {
		delete(c.items, key)
		return true
	}

	return false
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
			return true
		}

		c.Delete(key)
		return false
	}

	return false
}
