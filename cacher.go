// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

// Package cacher provides a uniform interface for different caching strategies.
package cacher

import "github.com/jelmersnoeck/cacher/memory"

// Cacher is the Caching interface that uniforms all the different strategies.
type Cacher interface {
	Add(key string, value []byte, ttl int64) error
	CompareAndReplace(token, key string, value []byte, ttl int64) error
	Set(key string, value []byte, ttl int64) error
	SetMulti(keys map[string][]byte, ttl int64) map[string]error
	Replace(key string, value []byte, ttl int64) error
	Increment(key string, initial, offset, ttl int64) error
	Decrement(key string, initial, offset, ttl int64) error
	Delete(key string) error
	DeleteMulti(keys []string) map[string]error
	Get(key string) ([]byte, string, error)
	GetMulti(keys []string) (map[string][]byte, map[string]string, map[string]error)
	Flush() error
	Touch(key string, ttl int64) error
}

// The default cache which will be used for package level functions.
var DefaultCache Cacher = memory.New(500)

// Add adds a value to the cache under the specified key with a given TTL.
func Add(key string, value []byte, ttl int64) error {
	return DefaultCache.Add(key, value, ttl)
}

// CompareAndReplace validates the token with the token in the store. If the
// tokens match, we will replace the value and return true. If it doesn't, we
// will not replace the value and return false.
func CompareAndReplace(token, key string, value []byte, ttl int64) error {
	return DefaultCache.CompareAndReplace(token, key, value, ttl)
}

// Set sets the value of an item, regardless of wether or not the value is
// already cached.
//
// ttl defines the number of seconds the value should be cached. If ttl is 0,
// the item will be cached infinitely.
func Set(key string, value []byte, ttl int64) error {
	return DefaultCache.Set(key, value, ttl)
}

// SetMulti sets multiple values for their respective keys. This is a shorthand
// to use `Set` multiple times.
func SetMulti(keys map[string][]byte, ttl int64) map[string]error {
	return DefaultCache.SetMulti(keys, ttl)
}

// Replace will update and only update the value of a cache key. If the key is
// not previously used, we will return false.
func Replace(key string, value []byte, ttl int64) error {
	return DefaultCache.Replace(key, value, ttl)
}

// Increment adds a value of offset to the initial value. If the initial value
// is already set, it will be added to the value currently stored in the cache.
func Increment(key string, initial, offset, ttl int64) error {
	return DefaultCache.Increment(key, initial, offset, ttl)
}

// Decrement subtracts a value of offset to the initial value. If the initial
// value is already set, it will be added to the value currently stored in the
// cache.
func Decrement(key string, initial, offset, ttl int64) error {
	return DefaultCache.Decrement(key, initial, offset, ttl)
}

// Delete will validate if the key actually is stored in the cache. If it is
// stored, it will remove the item from the cache. If it is not stored, it will
// return false.
func Delete(key string) error {
	return DefaultCache.Delete(key)
}

// DeleteMulti will delete multiple values at a time. It uses the `Delete`
// method internally to do so. It will return a map of results to see if the
// deletion is successful.
func DeleteMulti(keys []string) map[string]error {
	return DefaultCache.DeleteMulti(keys)
}

// Get gets the value out of the map associated with the provided key.
func Get(key string) ([]byte, string, error) {
	return DefaultCache.Get(key)
}

// GetMulti gets multiple values from the cache and returns them as a map. It
// uses `Get` internally to retrieve the data.
func GetMulti(keys []string) (map[string][]byte, map[string]string, map[string]error) {
	return DefaultCache.GetMulti(keys)
}

// Flush will remove all the items from the storage.
func Flush() error {
	return DefaultCache.Flush()
}

// Touch will update the key's ttl to the given ttl value without altering the
// value.
func Touch(key string, ttl int64) error {
	return DefaultCache.Touch(key, ttl)
}
