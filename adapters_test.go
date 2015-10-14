// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package cacher_test

import (
	"encoding/binary"
	"reflect"
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/jelmersnoeck/cacher"
	"github.com/jelmersnoeck/cacher/internal/numbers"
	"github.com/jelmersnoeck/cacher/internal/tests"
)

func TestAdd(t *testing.T) {
	for _, cache := range testDrivers() {
		if !cache.Add("key1", []byte("value1"), 0) {
			tests.FailMsg(t, cache, "Expecting `key1` to be added to the cache.")
		}

		if cache.Add("key1", []byte("value2"), 0) {
			tests.FailMsg(t, cache, "Expecting `key1` not to be added to the cache.")
		}

		tests.Compare(t, cache, "key1", "value1")

		cache.Flush()
	}
}

func TestSet(t *testing.T) {
	values := map[string][]byte{
		"key1": []byte("value"),
		"key2": numbers.Int64Bytes(2),
	}

	for _, cache := range testDrivers() {
		for key, value := range values {
			if !cache.Set(key, value, 0) {
				tests.FailMsg(t, cache, "Expecting `key1` to be `value`")
			}

			val, _ := cache.Get(key)
			if !reflect.DeepEqual(val, value) {
				tests.FailMsg(t, cache, "Value for key `"+key+"` does not match.")
			}
		}

		cache.Flush()
	}
}

func TestSetMulti(t *testing.T) {
	for _, cache := range testDrivers() {
		items := map[string][]byte{
			"item1": numbers.Int64Bytes(1),
			"item2": []byte("string"),
		}

		cache.SetMulti(items, 0)

		tests.Compare(t, cache, "item1", 1)
		tests.Compare(t, cache, "item2", "string")

		cache.Flush()
	}
}

func TestReplace(t *testing.T) {
	for _, cache := range testDrivers() {
		if cache.Replace("key1", []byte("value1"), 0) {
			tests.FailMsg(t, cache, "Key1 is not set yet, should not be able to replace.")
		}

		cache.Set("key1", []byte("value1"), 0)
		if !cache.Replace("key1", []byte("value1"), 0) {
			tests.FailMsg(t, cache, "Key1 has been set, should be able to replace.")
		}

		cache.Flush()
	}
}

func TestIncrement(t *testing.T) {
	for _, cache := range testDrivers() {
		cache.Increment("key1", 0, 1, 0)
		tests.Compare(t, cache, "key1", 0)

		cache.Increment("key1", 0, 1, 0)
		tests.Compare(t, cache, "key1", 1)

		cache.Set("string", []byte("value"), 0)
		if cache.Increment("string", 0, 1, 0) {
			tests.FailMsg(t, cache, "Can't increment a string value.")
		}

		if cache.Increment("key2", 0, 0, 0) {
			tests.FailMsg(t, cache, "Can't have an offset of <= 0")
		}

		if cache.Increment("key3", -1, 1, 0) {
			tests.FailMsg(t, cache, "Can't have an initial value of < 0")
		}

		cache.Flush()
	}
}

func TestDecrement(t *testing.T) {
	for _, cache := range testDrivers() {
		cache.Decrement("key1", 10, 1, 0)
		tests.Compare(t, cache, "key1", 10)

		cache.Decrement("key1", 10, 1, 0)
		tests.Compare(t, cache, "key1", 9)

		cache.Set("string", []byte("value"), 0)
		if cache.Decrement("string", 0, 1, 0) {
			tests.FailMsg(t, cache, "Can't decrement a string value.")
		}

		if cache.Decrement("key2", 0, 0, 0) {
			tests.FailMsg(t, cache, "Can't have an offset of <= 0")
		}

		if cache.Decrement("key3", -1, 1, 0) {
			tests.FailMsg(t, cache, "Can't have an initial value of < 0")
		}

		if cache.Decrement("key1", 10, 10, 0) {
			tests.FailMsg(t, cache, "Can't decrement below 0")
		}

		cache.Flush()
	}
}

func TestGet(t *testing.T) {
	for _, cache := range testDrivers() {
		cache.Set("key1", []byte("value1"), 0)
		tests.Compare(t, cache, "key1", "value1")

		if _, ok := cache.Get("key2"); ok {
			tests.FailMsg(t, cache, "Key2 is not present, ok should be false.")
		}

		cache.Flush()
	}
}

func TestGetMulti(t *testing.T) {
	for _, cache := range testDrivers() {
		items := map[string][]byte{
			"item1": numbers.Int64Bytes(1),
			"item2": []byte("string"),
		}

		cache.SetMulti(items, 0)

		var keys []string
		for k, _ := range items {
			keys = append(keys, k)
		}

		values := cache.GetMulti(keys)

		_, val := binary.Varint(values["item1"])
		if val != 1 {
			tests.FailMsg(t, cache, "Expected `item1` to equal `1`")
		}

		if string(values["item2"]) != "string" {
			tests.FailMsg(t, cache, "Expected `item2` to equal `string`")
		}

		cache.Flush()
	}
}

func TestDelete(t *testing.T) {
	for _, cache := range testDrivers() {
		cache.Set("key1", []byte("value1"), 0)
		tests.Compare(t, cache, "key1", "value1")

		cache.Delete("key1")

		if _, ok := cache.Get("key1"); ok {
			tests.FailMsg(t, cache, "`key1` should be deleted from the cache.")
		}

		cache.Flush()
	}
}

func TestDeleteMulti(t *testing.T) {
	for _, cache := range testDrivers() {
		items := map[string][]byte{
			"item1": numbers.Int64Bytes(1),
			"item2": []byte("string"),
		}

		cache.SetMulti(items, 0)
		cache.Set("key1", []byte("value1"), 0)

		var keys []string
		for k, _ := range items {
			keys = append(keys, k)
		}

		cache.DeleteMulti(keys)

		if _, ok := cache.Get("item1"); ok {
			tests.FailMsg(t, cache, "`item1` should be deleted from the cache.")
		}

		if _, ok := cache.Get("item2"); ok {
			tests.FailMsg(t, cache, "`item2` should be deleted from the cache.")
		}

		tests.Compare(t, cache, "key1", "value1")

		cache.Flush()
	}
}

func TestFlush(t *testing.T) {
	for _, cache := range testDrivers() {
		cache.Set("key1", []byte("value1"), 0)
		tests.Compare(t, cache, "key1", "value1")

		if !cache.Flush() {
			tests.FailMsg(t, cache, "Cache should be able to flush")
		}

		if v, _ := cache.Get("key1"); v != nil {
			tests.FailMsg(t, cache, "Expecting `key1` to be nil")
		}

		cache.Flush()
	}
}

func testDrivers() []cacher.Cacher {
	drivers := make([]cacher.Cacher, 0)
	drivers = append(drivers, cacher.NewMemoryCache(0))

	c, _ := redis.Dial("tcp", ":6379")
	redisCache := cacher.NewRedisCache(c)
	drivers = append(drivers, redisCache)

	return drivers
}