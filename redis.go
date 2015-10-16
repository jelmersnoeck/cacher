// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package cacher

import (
	"github.com/garyburd/redigo/redis"
	"github.com/jelmersnoeck/cacher/internal/numbers"
)

// RedisCache is a caching implementation that stores the data in memory. The
// cache will be emptied when the application has run.
type RedisCache struct {
	client redis.Conn
}

// NewRedisCache creates a new instance of RedisCache and initiates the
// storage map.
func NewRedisCache(client redis.Conn) *RedisCache {
	cache := new(RedisCache)
	cache.client = client

	return cache
}

// Add an item to the cache. If the item is already cached, the value won't be
// overwritten.
//
// ttl defines the number of seconds the value should be cached. If ttl is 0,
// the item will be cached infinitely.
func (c *RedisCache) Add(key string, value []byte, ttl int64) bool {
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
func (c *RedisCache) Set(key string, value []byte, ttl int64) bool {
	var err error

	if ttl > 0 {
		_, err = c.client.Do("SETEX", key, ttl, value)
	} else if ttl == 0 {
		_, err = c.client.Do("SET", key, value)
	} else {
		return c.Delete(key)
	}

	return err == nil
}

// SetMulti sets multiple values for their respective keys. This is a shorthand
// to use `Set` multiple times.
func (c *RedisCache) SetMulti(items map[string][]byte, ttl int64) map[string]bool {
	results := make(map[string]bool)

	c.client.Do("MULTI")
	for key, value := range items {
		results[key] = c.Set(key, value, ttl)
	}
	c.client.Do("EXEC")

	return results
}

// Replace will update and only update the value of a cache key. If the key is
// not previously used, we will return false.
func (c *RedisCache) Replace(key string, value []byte, ttl int64) bool {
	if !c.exists(key) {
		return false
	}

	return c.Set(key, value, ttl)
}

// Get gets the value out of the map associated with the provided key.
func (c *RedisCache) Get(key string) ([]byte, bool) {
	value, _ := c.client.Do("GET", key)

	if value == nil {
		return []byte{}, false
	}

	val, ok := value.([]byte)

	if !ok {
		return nil, false
	}

	return val, true
}

// GetMulti gets multiple values from the cache and returns them as a map. It
// uses `Get` internally to retrieve the data.
func (c *RedisCache) GetMulti(keys []string) map[string][]byte {
	cValues, err := c.client.Do("MGET", keyArgs(keys)...)
	items := make(map[string][]byte)

	if err == nil {
		values := cValues.([]interface{})
		for i, val := range values {
			items[keys[i]] = val.([]byte)
		}
	}

	return items
}

// Increment adds a value of offset to the initial value. If the initial value
// is already set, it will be added to the value currently stored in the cache.
func (c *RedisCache) Increment(key string, initial, offset, ttl int64) bool {
	if initial < 0 || offset <= 0 {
		return false
	}

	return c.incrementOffset(key, initial, offset, ttl)
}

// Decrement subtracts a value of offset to the initial value. If the initial
// value is already set, it will be added to the value currently stored in the
// cache.
func (c *RedisCache) Decrement(key string, initial, offset, ttl int64) bool {
	if initial < 0 || offset <= 0 {
		return false
	}

	return c.incrementOffset(key, initial, offset*-1, ttl)
}

// Flush will remove all the items from the hash.
func (c *RedisCache) Flush() bool {
	_, err := c.client.Do("FLUSHDB")

	return err == nil
}

// Delete will validate if the key actually is stored in the cache. If it is
// stored, it will remove the item from the cache. If it is not stored, it will
// return false.
func (c *RedisCache) Delete(key string) bool {
	_, err := c.client.Do("DEL", key)

	if err != nil {
		return false
	}

	return true
}

// DeleteMulti will delete multiple values at a time. It uses the `Delete`
// method internally to do so. It will return a map of results to see if the
// deletion is successful.
func (c *RedisCache) DeleteMulti(keys []string) map[string]bool {
	items := c.GetMulti(keys)
	c.client.Do("DEL", keyArgs(keys)...)

	results := make(map[string]bool)
	for _, v := range keys {
		_, results[v] = items[v]
	}

	return results
}

// incrementOffset is a common incrementor method used between Increment and
// Decrement. If the key isn't set before, we will set the initial value. If
// there is a value present, we will add the given offset to that value and
// update the value with the new TTL.
func (c *RedisCache) incrementOffset(key string, initial, offset, ttl int64) bool {
	c.client.Do("WATCH", key)

	if !c.exists(key) {
		c.client.Do("MULTI")
		defer c.client.Do("EXEC")
		return c.Set(key, numbers.Int64Bytes(initial), ttl)
	}

	getValue, _ := c.Get(key)
	val, ok := numbers.BytesInt64(getValue)

	if !ok {
		return false
	}

	c.client.Do("MULTI")
	defer c.client.Do("EXEC")

	val += offset
	if val < 0 {
		return false
	}

	return c.Set(key, numbers.Int64Bytes(val), ttl)
}

func (c *RedisCache) exists(key string) bool {
	val, _ := c.client.Do("EXISTS", key)

	if val.(int64) == 1 {
		return true
	}

	return false
}

func keyArgs(keys []string) []interface{} {
	var args []interface{}
	for _, key := range keys {
		args = append(args, key)
	}

	return args
}
