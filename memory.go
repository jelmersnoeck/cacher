// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package cacher

// MemoryCache is a caching implementation that stores the data in memory. The
// cache will be emptied when the application has run.
type MemoryCache struct {
	items map[string]interface{}
}

// NewMemoryCache creates a new instance of MemoryCache and initiates the
// storage map.
func NewMemoryCache() *MemoryCache {
	cache := new(MemoryCache)
	cache.items = make(map[string]interface{})

	return cache
}

// Set adds the value to the specified key in the map.
func (c *MemoryCache) Set(key string, value interface{}) bool {
	c.items[key] = value
	return true
}

// Get gets the value out of the map associated with the provided key.
func (c *MemoryCache) Get(key string) interface{} {
	return c.items[key]
}
