// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package memory

import (
	"runtime"
	"time"

	"github.com/jelmersnoeck/cacher/internal/encoding"
)

type cachedItem struct {
	value  []byte
	expiry time.Time
	expire bool
	token  string
}

// Cache is a caching implementation that stores the data in memory. The
// cache will be emptied when the application has run.
type Cache struct {
	items map[string]*cachedItem
	keys  []string
	limit uintptr
	size  uintptr
}

// New creates a new instance of Cache and initiates the storage map.
func New(limit uintptr) *Cache {
	cache := new(Cache)
	cache.items = make(map[string]*cachedItem)
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
func (c *Cache) Add(key string, value []byte, ttl int64) bool {
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
func (c *Cache) Set(key string, value []byte, ttl int64) bool {
	expiry := time.Now().Add(time.Duration(ttl) * time.Second)

	var expire bool
	if ttl != 0 {
		expire = true
	}

	c.items[key] = &cachedItem{value, expiry, expire, encoding.Md5Sum(value)}
	c.keys = append(c.keys, key)
	c.size += uintptr(len(value)) // TODO: if already exists, don't add this all
	c.lru(key)
	c.evict()
	return true
}

// SetMulti sets multiple values for their respective keys. This is a shorthand
// to use `Set` multiple times.
func (c *Cache) SetMulti(items map[string][]byte, ttl int64) map[string]bool {
	results := make(map[string]bool)
	for key, value := range items {
		results[key] = c.Set(key, value, ttl)
	}

	return results
}

// CompareAndReplace validates the token with the token in the store. If the
// tokens match, we will replace the value and return true. If it doesn't, we
// will not replace the value and return false.
func (c *Cache) CompareAndReplace(token, key string, value []byte, ttl int64) bool {
	if !c.exists(key) {
		return false
	}

	if c.items[key].token != token {
		return false
	}

	return c.Set(key, value, ttl)
}

// Replace will update and only update the value of a cache key. If the key is
// not previously used, we will return false.
func (c *Cache) Replace(key string, value []byte, ttl int64) bool {
	if !c.exists(key) {
		return false
	}

	return c.Set(key, value, ttl)
}

// Get gets the value out of the map associated with the provided key.
func (c *Cache) Get(key string) ([]byte, string, bool) {
	if c.exists(key) {
		return c.items[key].value, c.items[key].token, true
	}

	return nil, "", false
}

// GetMulti gets multiple values from the cache and returns them as a map. It
// uses `Get` internally to retrieve the data.
func (c *Cache) GetMulti(keys []string) (map[string][]byte, map[string]string, map[string]bool) {
	items := make(map[string][]byte)
	bools := make(map[string]bool)
	tokens := make(map[string]string)

	for _, k := range keys {
		items[k], tokens[k], bools[k] = c.Get(k)
	}

	return items, tokens, bools
}

// Increment adds a value of offset to the initial value. If the initial value
// is already set, it will be added to the value currently stored in the cache.
func (c *Cache) Increment(key string, initial, offset, ttl int64) bool {
	if initial < 0 || offset <= 0 {
		return false
	}

	return c.incrementOffset(key, initial, offset, ttl)
}

// Decrement subtracts a value of offset to the initial value. If the initial
// value is already set, it will be added to the value currently stored in the
// cache.
func (c *Cache) Decrement(key string, initial, offset, ttl int64) bool {
	if initial < 0 || offset <= 0 {
		return false
	}

	return c.incrementOffset(key, initial, offset*-1, ttl)
}

// Flush will remove all the items from the hash.
func (c *Cache) Flush() bool {
	c.items = make(map[string]*cachedItem)
	c.size = 0
	return true
}

// Delete will validate if the key actually is stored in the cache. If it is
// stored, it will remove the item from the cache. If it is not stored, it will
// return false.
func (c *Cache) Delete(key string) bool {
	for i, v := range c.keys {
		if v == key {
			return c.removeAt(i)
		}
	}

	return false
}

// DeleteMulti will delete multiple values at a time. It uses the `Delete`
// method internally to do so. It will return a map of results to see if the
// deletion is successful.
func (c *Cache) DeleteMulti(keys []string) map[string]bool {
	results := make(map[string]bool)

	for _, key := range keys {
		results[key] = c.Delete(key)
	}

	return results
}

// Touch will update the key's ttl to the given ttl value without altering the
// value.
func (c *Cache) Touch(key string, ttl int64) bool {
	if !c.exists(key) {
		return false
	}

	if ttl < 0 {
		return c.Delete(key)
	}

	c.items[key].expiry = time.Now().Add(time.Duration(ttl) * time.Second)
	return true
}

// removeAt will remove a specific indexed value from our cache.
func (c *Cache) removeAt(index int) bool {
	key := c.keys[index]
	c.keys = append(c.keys[:index], c.keys[index+1:]...)
	c.size -= uintptr(len(c.items[key].value))
	delete(c.items, key)

	return true
}

// incrementOffset is a common incrementor method used between Increment and
// Decrement. If the key isn't set before, we will set the initial value. If
// there is a value present, we will add the given offset to that value and
// update the value with the new TTL.
func (c *Cache) incrementOffset(key string, initial, offset, ttl int64) bool {
	if !c.exists(key) {
		return c.Set(key, encoding.Int64Bytes(initial), ttl)
	}

	val, ok := encoding.BytesInt64(c.items[key].value)

	if !ok {
		return false
	}

	val += offset
	if val < 0 {
		return false
	}

	return c.Set(key, encoding.Int64Bytes(val), ttl)
}

// exists checks if a key is stored in the cache.
//
// If the key is stored in the cache, but the expiry date has passed, we will
// remove the item from the cache and return false. If the expiry has not passed
// yet, it will return false.
func (c *Cache) exists(key string) bool {
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
func (c *Cache) evict() {
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
func (c *Cache) lru(key string) {
	for i, v := range c.keys {
		if v == key {
			c.keys = append(c.keys[:i], c.keys[i+1:]...)
			c.keys = append(c.keys, key)
		}
	}
}
