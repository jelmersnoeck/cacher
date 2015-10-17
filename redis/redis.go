// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package redis

import (
	"github.com/garyburd/redigo/redis"
	"github.com/jelmersnoeck/cacher/internal/encoding"
)

// Cache is an instance that stores a Redis client that will be used to
// communicate with the Redis server.
type Cache struct {
	client redis.Conn
}

// New creates a new instance of Cache.
func New(client redis.Conn) *Cache {
	cache := new(Cache)
	cache.client = client

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
func (c *Cache) SetMulti(items map[string][]byte, ttl int64) map[string]bool {
	results := make(map[string]bool)

	c.client.Do("MULTI")
	for key, value := range items {
		results[key] = c.Set(key, value, ttl)
	}
	c.client.Do("EXEC")

	return results
}

// CompareAndReplace validates the token with the token in the store. If the
// tokens match, we will replace the value and return true. If it doesn't, we
// will not replace the value and return false.
func (c *Cache) CompareAndReplace(token, key string, value []byte, ttl int64) bool {
	c.client.Do("WATCH", key)
	defer c.client.Do("UNWATCH")

	if !c.exists(key) {
		return false
	}

	_, storedToken, _ := c.Get(key)
	if token != storedToken {
		return false
	}

	// We're watching the key, by using MULTI the transaction will fail if the key
	// changes in the meantime.
	c.client.Do("MULTI")
	c.Set(key, value, ttl)
	rValue, _ := c.client.Do("EXEC")

	for _, v := range rValue.([]interface{}) {
		if v.(string) != "OK" {
			return false
		}
	}

	return true
}

// Replace will update and only update the value of a cache key. If the key is
// not previously used, we will return false.
func (c *Cache) Replace(key string, value []byte, ttl int64) bool {
	c.client.Do("WATCH", key)
	defer c.client.Do("UNWATCH")

	if !c.exists(key) {
		return false
	}

	// We're watching the key, so we can use a transaction to set the value. If
	// the key changes in the meantime, it'll fail.
	c.client.Do("MULTI")
	c.Set(key, value, ttl)
	vals, err := c.client.Do("EXEC")

	if err != nil {
		return false
	}

	for _, v := range vals.([]interface{}) {
		if v.(string) != "OK" {
			return false
		}
	}

	return true
}

// Get gets the value out of the map associated with the provided key.
func (c *Cache) Get(key string) ([]byte, string, bool) {
	value, _ := c.client.Do("GET", key)

	if value == nil {
		return []byte{}, "", false
	}

	val, ok := value.([]byte)

	if !ok {
		return nil, "", false
	}

	return val, encoding.Md5Sum(val), true
}

// GetMulti gets multiple values from the cache and returns them as a map. It
// uses `Get` internally to retrieve the data.
func (c *Cache) GetMulti(keys []string) (map[string][]byte, map[string]string, map[string]bool) {
	cValues, err := c.client.Do("MGET", keyArgs(keys)...)
	items := make(map[string][]byte)
	bools := make(map[string]bool)
	tokens := make(map[string]string)

	for _, v := range keys {
		bools[v] = false
	}

	if err == nil {
		values := cValues.([]interface{})
		for i, val := range values {
			items[keys[i]] = val.([]byte)
			tokens[keys[i]] = encoding.Md5Sum(items[keys[i]])
			bools[keys[i]] = true
		}
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
	_, err := c.client.Do("FLUSHDB")

	return err == nil
}

// Delete will validate if the key actually is stored in the cache. If it is
// stored, it will remove the item from the cache. If it is not stored, it will
// return false.
func (c *Cache) Delete(key string) bool {
	_, err := c.client.Do("DEL", key)

	if err != nil {
		return false
	}

	return true
}

// DeleteMulti will delete multiple values at a time. It uses the `Delete`
// method internally to do so. It will return a map of results to see if the
// deletion is successful.
func (c *Cache) DeleteMulti(keys []string) map[string]bool {
	items, _, _ := c.GetMulti(keys)
	c.client.Do("DEL", keyArgs(keys)...)

	// DEL will only return false if the key is not present. To get a map of bools
	// to return, we can go over the items that are in the store (before we've
	// deleted them) and see which of the specified keys to delete are present in
	// the list of items.
	results := make(map[string]bool)
	for _, key := range keys {
		_, results[key] = items[key]
	}

	return results
}

// incrementOffset is a common incrementor method used between Increment and
// Decrement. If the key isn't set before, we will set the initial value. If
// there is a value present, we will add the given offset to that value and
// update the value with the new TTL.
func (c *Cache) incrementOffset(key string, initial, offset, ttl int64) bool {
	c.client.Do("WATCH", key)

	if !c.exists(key) {
		c.client.Do("MULTI")
		defer c.client.Do("EXEC")
		return c.Set(key, encoding.Int64Bytes(initial), ttl)
	}

	getValue, _, _ := c.Get(key)
	val, ok := encoding.BytesInt64(getValue)

	if !ok {
		return false
	}

	// We are watching our key. With using a transaction, we can check that this
	// increment doesn't inflect with another concurrent request that might
	// happen.
	c.client.Do("MULTI")
	defer c.client.Do("EXEC")

	val += offset
	if val < 0 {
		return false
	}

	return c.Set(key, encoding.Int64Bytes(val), ttl)
}

func (c *Cache) exists(key string) bool {
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
