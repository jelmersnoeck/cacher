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
	"github.com/jelmersnoeck/cacher/internal/encoding"
	"github.com/jelmersnoeck/cacher/internal/tests"
	"github.com/jelmersnoeck/cacher/memory"
	rcache "github.com/jelmersnoeck/cacher/redis"
)

func TestAdd(t *testing.T) {
	for _, cache := range testDrivers() {
		if err := cache.Add("key1", []byte("value1"), 0); err != nil {
			tests.FailMsg(t, cache, "Expecting `key1` to be added to the cache.")
		}

		if err := cache.Add("key1", []byte("value2"), 0); err == nil {
			tests.FailMsg(t, cache, "Expecting `key1` not to be added to the cache.")
		}

		tests.Compare(t, cache, "key1", "value1")
	}
}

func TestSet(t *testing.T) {
	values := map[string][]byte{
		"key1": []byte("value"),
		"key2": encoding.Int64Bytes(2),
	}

	for _, cache := range testDrivers() {
		for key, value := range values {
			if err := cache.Set(key, value, 0); err != nil {
				tests.FailMsg(t, cache, "Expecting `key1` to be `value`")
			}

			val, _, _ := cache.Get(key)
			if !reflect.DeepEqual(val, value) {
				tests.FailMsg(t, cache, "Value for key `"+key+"` does not match.")
			}
		}

		cache.Set("key1", []byte("value"), -1)
		_, _, err := cache.Get("key1")

		if err == nil {
			tests.FailMsg(t, cache, "key1 should be deleted with negative value")
		}

	}
}

func TestSetMulti(t *testing.T) {
	for _, cache := range testDrivers() {
		items := map[string][]byte{
			"item1": encoding.Int64Bytes(1),
			"item2": []byte("string"),
		}

		cache.SetMulti(items, 0)

		tests.Compare(t, cache, "item1", 1)
		tests.Compare(t, cache, "item2", "string")
	}
}

func TestCompareAndReplace(t *testing.T) {
	for _, cache := range testDrivers() {
		cache.Set("key1", []byte("CompareAndReplace"), 0)
		val1, token1, _ := cache.Get("key1")
		if string(val1) != "CompareAndReplace" {
			tests.FailMsg(t, cache, "`key1` should equal `CompareAndReplace`")
		}

		err := cache.CompareAndReplace(token1, "key1", []byte("ReplacementValue"), 0)
		if err != nil {
			tests.FailMsg(t, cache, "CompareAndReplace should be executed.")
		}
		val2, token2, _ := cache.Get("key1")
		if string(val2) != "ReplacementValue" {
			tests.FailMsg(t, cache, "`key1` should equal `ReplacementValue`")
		}

		err = cache.CompareAndReplace(token2+"WRONG", "key1", []byte("WrongValue"), 0)
		if err == nil {
			tests.FailMsg(t, cache, "WrongValue should not be set.")
		}
		val3, _, _ := cache.Get("key1")
		if string(val3) != "ReplacementValue" {
			tests.FailMsg(t, cache, "`key1` should equal `ReplacementValue`")
		}

	}
}

func TestReplace(t *testing.T) {
	for _, cache := range testDrivers() {
		if err := cache.Replace("key1", []byte("value1"), 0); err == nil {
			tests.FailMsg(t, cache, "Key1 is not set yet, should not be able to replace.")
		}

		cache.Set("key1", []byte("value1"), 0)
		if err := cache.Replace("key1", []byte("value1"), 0); err != nil {
			tests.FailMsg(t, cache, "Key1 has been set, should be able to replace.")
		}

	}
}

func TestCache_Increment(t *testing.T) {
	for _, cache := range testDrivers() {
		cache.Increment("key1", 0, 1, 0)
		cache.Increment("key1", 0, 1, 0)
		v, _, _ := cache.Get("key1")

		num, _ := encoding.BytesInt64(v)
		if num != 1 {
			tests.FailMsg(t, cache, "Expected the value to be 1, got %d", num)
		}

		cache.Set("key2", []byte("string value, not incrementable"), 0)
		err := cache.Increment("key2", 0, 5, 0)
		if err == nil {
			tests.FailMsg(t, cache, "Expected the error not to be nil")
		}
	}
}

func TestIncrement(t *testing.T) {
	for _, cache := range testDrivers() {
		cache.Increment("key1", 0, 1, 0)
		tests.Compare(t, cache, "key1", 0)

		cache.Increment("key1", 0, 1, 0)
		tests.Compare(t, cache, "key1", 1)

		cache.Set("string", []byte("value"), 0)
		if err := cache.Increment("string", 0, 1, 0); err == nil {
			tests.FailMsg(t, cache, "Can't increment a string value.")
		}

		if err := cache.Increment("key2", 0, 0, 0); err == nil {
			tests.FailMsg(t, cache, "Can't have an offset of <= 0")
		}

		if err := cache.Increment("key3", -1, 1, 0); err == nil {
			tests.FailMsg(t, cache, "Can't have an initial value of < 0")
		}

	}
}

func TestCache_Decrement(t *testing.T) {
	for _, cache := range testDrivers() {
		cache.Decrement("key1", 10, 1, 0)
		cache.Decrement("key1", 10, 3, 0)
		v, _, _ := cache.Get("key1")
		num, _ := encoding.BytesInt64(v)

		if num != 7 {
			tests.FailMsg(t, cache, "Expected value to be 7, got %d", num)
		}

		cache.Set("key2", []byte("string value, not decrementable"), 0)
		err := cache.Decrement("key2", 0, 5, 0)

		if err == nil {
			tests.FailMsg(t, cache, "Expected error not to be nil")
		}
	}
}

func TestDecrement(t *testing.T) {
	for _, cache := range testDrivers() {
		cache.Decrement("key1", 10, 1, 0)
		tests.Compare(t, cache, "key1", 10)

		cache.Decrement("key1", 10, 1, 0)
		tests.Compare(t, cache, "key1", 9)

		cache.Set("string", []byte("value"), 0)
		if err := cache.Decrement("string", 0, 1, 0); err == nil {
			tests.FailMsg(t, cache, "Can't decrement a string value.")
		}

		if err := cache.Decrement("key2", 0, 0, 0); err == nil {
			tests.FailMsg(t, cache, "Can't have an offset of <= 0")
		}

		if err := cache.Decrement("key3", -1, 1, 0); err == nil {
			tests.FailMsg(t, cache, "Can't have an initial value of < 0")
		}

		if err := cache.Decrement("key1", 10, 10, 0); err == nil {
			tests.FailMsg(t, cache, "Can't decrement below 0")
		}

	}
}

func TestGet(t *testing.T) {
	for _, cache := range testDrivers() {
		cache.Set("key1", []byte("value1"), 0)
		tests.Compare(t, cache, "key1", "value1")

		if _, _, err := cache.Get("key2"); err == nil {
			tests.FailMsg(t, cache, "Key2 is not present, err should not be nil.")
		}

	}
}

func TestGetToken(t *testing.T) {
	for _, cache := range testDrivers() {
		cache.Set("key1", []byte("value1"), 0)
		_, token1, _ := cache.Get("key1")

		cache.Set("key1", []byte("value2"), 0)
		_, token2, _ := cache.Get("key1")

		if token1 == token2 {
			tests.FailMsg(t, cache, "token1 should not equal token2.")
		}

	}
}

func TestGetMulti(t *testing.T) {
	for _, cache := range testDrivers() {
		items := map[string][]byte{
			"item1": encoding.Int64Bytes(1),
			"item2": []byte("string"),
		}

		cache.SetMulti(items, 0)

		var keys []string
		for k := range items {
			keys = append(keys, k)
		}

		values, tokens, errs := cache.GetMulti(keys)

		_, val := binary.Varint(values["item1"])
		if val != 1 {
			tests.FailMsg(t, cache, "Expected `item1` to equal `1`")
		}

		if err, ok := errs["item1"]; !ok || err != nil {
			tests.FailMsg(t, cache, "Expected `item1` to be ok.")
		}

		if tokens["item1"] == "" {
			tests.FailMsg(t, cache, "Expected `item1` to have a valid token.")
		}

		if string(values["item2"]) != "string" {
			tests.FailMsg(t, cache, "Expected `item2` to equal `string`")
		}
	}
}

func TestDelete(t *testing.T) {
	for _, cache := range testDrivers() {
		cache.Set("key1", []byte("value1"), 0)
		tests.Compare(t, cache, "key1", "value1")

		cache.Delete("key1")

		if _, _, err := cache.Get("key1"); err == nil {
			tests.FailMsg(t, cache, "`key1` should be deleted from the cache.")
		}
	}
}

func TestDeleteMulti(t *testing.T) {
	for _, cache := range testDrivers() {
		items := map[string][]byte{
			"item1": encoding.Int64Bytes(1),
			"item2": []byte("string"),
		}

		cache.SetMulti(items, 0)
		cache.Set("key1", []byte("value1"), 0)

		var keys []string
		for k := range items {
			keys = append(keys, k)
		}

		cache.DeleteMulti(keys)

		if _, _, err := cache.Get("item1"); err == nil {
			tests.FailMsg(t, cache, "`item1` should be deleted from the cache.")
		}

		if _, _, err := cache.Get("item2"); err == nil {
			tests.FailMsg(t, cache, "`item2` should be deleted from the cache.")
		}

		tests.Compare(t, cache, "key1", "value1")
	}
}

func TestFlush(t *testing.T) {
	for _, cache := range testDrivers() {
		cache.Set("key1", []byte("value1"), 0)
		tests.Compare(t, cache, "key1", "value1")

		if err := cache.Flush(); err != nil {
			tests.FailMsg(t, cache, "Cache should be able to flush")
		}

		if _, _, err := cache.Get("key1"); err == nil {
			tests.FailMsg(t, cache, "Expecting `key1` to be nil")
		}
	}
}

func TestTouch(t *testing.T) {
	for _, cache := range testDrivers() {
		if err := cache.Touch("key1", 5); err == nil {
			tests.FailMsg(t, cache, "Can't touch a non-existing key.")
		}

		cache.Set("key1", []byte("Hello world"), 0)
		if err := cache.Touch("key1", 5); err != nil {
			tests.FailMsg(t, cache, "Should be able to touch existing key.")
		}
	}
}

func testDrivers() []cacher.Cacher {
	var drivers []cacher.Cacher

	memoryCache := memory.New(0)
	memoryCache.Flush()
	drivers = append(drivers, memoryCache)

	c, _ := redis.Dial("tcp", ":6379")
	redisCache := rcache.New(c)
	redisCache.Flush()
	drivers = append(drivers, redisCache)

	return drivers
}
