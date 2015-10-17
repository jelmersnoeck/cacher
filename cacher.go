// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

// Package cacher provides a uniform interface for different caching strategies.
package cacher

// Cacher is the Caching interface that uniforms all the different strategies.
type Cacher interface {
	Add(key string, value []byte, ttl int64) bool
	CompareAndReplace(token, key string, value []byte, ttl int64) bool
	Set(key string, value []byte, ttl int64) bool
	SetMulti(keys map[string][]byte, ttl int64) map[string]bool
	Replace(key string, value []byte, ttl int64) bool
	Increment(key string, initial, offset, ttl int64) bool
	Decrement(key string, initial, offset, ttl int64) bool
	Delete(key string) bool
	DeleteMulti(keys []string) map[string]bool
	Get(key string) ([]byte, string, bool)
	GetMulti(keys []string) (map[string][]byte, map[string]string, map[string]bool)
	Flush() bool
	Touch(key string, ttl int64) bool
}
