// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package redis_test

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
	"github.com/jelmersnoeck/cacher"
	"github.com/jelmersnoeck/cacher/internal/encoding"
	rcache "github.com/jelmersnoeck/cacher/redis"
)

func ExampleCache_Add() {
	cache := getCache()

	ok1 := cache.Add("key1", []byte("value1"), 0)
	val1, _, _ := cache.Get("key1")
	fmt.Println(string(val1), ok1)

	ok2 := cache.Add("key1", []byte("value2"), 0)
	val2, _, _ := cache.Get("key1")
	fmt.Println(string(val2), ok2)

	// Output:
	// value1 true
	// value1 false
}

func ExampleCache_Set() {
	cache := getCache()

	ok1 := cache.Set("key1", []byte("value1"), 0)
	val1, _, _ := cache.Get("key1")
	fmt.Println(string(val1), ok1)

	ok2 := cache.Set("key1", []byte("value2"), 0)
	val2, _, _ := cache.Get("key1")
	fmt.Println(string(val2), ok2)

	// Output:
	// value1 true
	// value2 true
}

func ExampleCache_SetMulti() {
	cache := getCache()

	multi := map[string][]byte{
		"key1": []byte("value1"),
		"key2": []byte("value2"),
	}

	cache.SetMulti(multi, 0)
	val1, _, _ := cache.Get("key1")
	val2, _, _ := cache.Get("key2")

	fmt.Println(string(val1))
	fmt.Println(string(val2))

	// Output:
	// value1
	// value2
}

func ExampleCache_CompareAndReplace() {
	cache := getCache()

	cache.Set("key1", []byte("hello world"), 0)
	_, token, _ := cache.Get("key1")

	var ok bool
	var val []byte
	ok = cache.CompareAndReplace(token+"FALSE", "key1", []byte("replacement1"), 0)
	val, _, _ = cache.Get("key1")
	fmt.Println(ok, string(val))

	ok = cache.CompareAndReplace(token, "key1", []byte("replacement2"), 0)
	val, _, _ = cache.Get("key1")
	fmt.Println(ok, string(val))

	// Output:
	// false hello world
	// true replacement2
}

func ExampleCache_Replace() {
	cache := getCache()
	var ok bool

	ok = cache.Replace("key1", []byte("replacement"), 0)
	fmt.Println(ok)

	cache.Set("key1", []byte("value1"), 0)
	ok = cache.Replace("key1", []byte("replacement"), 0)
	fmt.Println(ok)

	// Output:
	// false
	// true
}

func ExampleCache_Get() {
	cache := getCache()
	var value []byte
	var token string
	var ok bool

	value, token, ok = cache.Get("non-existing")
	fmt.Println(string(value), token, ok)

	cache.Set("key1", []byte("Hello world!"), 0)
	value, token, ok = cache.Get("key1")
	fmt.Println(string(value), token, ok)

	// Output:
	// false
	// Hello world! 86fb269d190d2c85f6e0468ceca42a20 true
}

func ExampleCache_GetMulti() {
	cache := getCache()

	multi := map[string][]byte{
		"key1": []byte("value1"),
		"key2": []byte("value2"),
	}

	cache.SetMulti(multi, 0)

	keys := []string{"key1", "key2"}

	values, tokens, bools := cache.GetMulti(keys)
	fmt.Println(values["key1"], values["key2"])
	fmt.Println(tokens["key1"], tokens["key2"])
	fmt.Println(bools["key1"], bools["key2"])

	// Output:
	// [118 97 108 117 101 49] [118 97 108 117 101 50]
	// 9946687e5fa0dab5993ededddb398d2e f066ce9385512ee02afc6e14d627e9f2
	// true true
}

func ExampleCache_Increment() {
	cache := getCache()

	cache.Increment("key1", 0, 1, 0)
	cache.Increment("key1", 0, 1, 0)
	v, _, _ := cache.Get("key1")

	num, _ := encoding.BytesInt64(v)
	fmt.Println(num)

	cache.Set("key2", []byte("string value, not incrementable"), 0)
	ok := cache.Increment("key2", 0, 5, 0)
	v2, _, _ := cache.Get("key2")
	fmt.Println(ok, string(v2))

	// Output:
	// 1
	// false string value, not incrementable
}

func ExampleCache_Decrement() {
	cache := getCache()

	cache.Decrement("key1", 10, 1, 0)
	cache.Decrement("key1", 10, 3, 0)
	v, _, _ := cache.Get("key1")

	num, _ := encoding.BytesInt64(v)
	fmt.Println(num)

	cache.Set("key2", []byte("string value, not decrementable"), 0)
	ok := cache.Decrement("key2", 0, 5, 0)
	v2, _, _ := cache.Get("key2")
	fmt.Println(ok, string(v2))

	// Output:
	// 7
	// false string value, not decrementable
}

func ExampleCache_Flush() {
	cache := getCache()
	var bools map[string]bool

	keys := []string{"key1", "key2"}
	multi := map[string][]byte{
		"key1": []byte("value1"),
		"key2": []byte("value2"),
	}
	cache.SetMulti(multi, 0)

	_, _, bools = cache.GetMulti(keys)
	fmt.Println(bools["key1"], bools["key2"])

	cache.Flush()
	_, _, bools = cache.GetMulti(keys)
	fmt.Println(bools["key1"], bools["key2"])

	// Output:
	// true true
	// false false
}

func ExampleCache_Delete() {
	cache := getCache()
	var ok bool

	cache.Set("key1", []byte("value1"), 0)
	ok = cache.Delete("key1")
	fmt.Println(ok)

	ok = cache.Delete("non-existing")
	fmt.Println(ok)

	// Output:
	// true
	// false
}

func ExampleCache_DeleteMulti() {
	cache := getCache()

	keys := []string{"key1", "key2", "non-existing"}

	multi := map[string][]byte{
		"key1": []byte("value1"),
		"key2": []byte("value2"),
	}
	cache.SetMulti(multi, 0)

	oks := cache.DeleteMulti(keys)
	fmt.Println(oks["key1"])
	fmt.Println(oks["key2"])
	fmt.Println(oks["non-existing"])

	// Output:
	// true
	// true
	// false
}

func ExampleCache_Touch() {
	cache := getCache()

	cache.Set("key1", []byte("value1"), 1)
	ok := cache.Touch("key1", 5)
	fmt.Println(ok)

	// Output:
	// true
}

func getCache() cacher.Cacher {
	c, _ := redis.Dial("tcp", ":6379")
	cache := rcache.New(c)
	cache.Flush()

	return cache
}
