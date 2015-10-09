// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package cacher

import "time"

type cachedItem struct {
	value  interface{}
	expiry time.Time
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
// ttl defines the number of seconds the value should be cached.
func (c *MemoryCache) Add(key string, value interface{}, ttl int) bool {
	_, exists := c.items[key]
	if exists {
		return false
	}

	return c.Set(key, value, ttl)
}

// Set sets the value of an item, regardless of wether or not the value is
// already cached.
//
// ttl defines the number of seconds the value should be cached.
func (c *MemoryCache) Set(key string, value interface{}, ttl int) bool {
	expiry := time.Now().Add(time.Duration(ttl) * time.Second)
	c.items[key] = cachedItem{value, expiry}
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
