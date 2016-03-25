// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package cacher_test

import (
	"fmt"

	"github.com/jelmersnoeck/cacher"
	"github.com/jelmersnoeck/cacher/internal/encoding"
)

func ExampleCache_Add() {
	cacher.Flush()
	ok1 := cacher.Add("key1", []byte("value1"), 0)
	val1, _, _ := cacher.Get("key1")
	fmt.Println(string(val1), ok1)

	ok2 := cacher.Add("key1", []byte("value2"), 0)
	val2, _, _ := cacher.Get("key1")
	fmt.Println(string(val2), ok2)

	// Output:
	// value1 <nil>
	// value1 Key `key1` already exists.
}

func ExampleCache_Set() {
	cacher.Flush()
	ok1 := cacher.Set("key1", []byte("value1"), 0)
	val1, _, _ := cacher.Get("key1")
	fmt.Println(string(val1), ok1)

	ok2 := cacher.Set("key1", []byte("value2"), 0)
	val2, _, _ := cacher.Get("key1")
	fmt.Println(string(val2), ok2)

	// Output:
	// value1 <nil>
	// value2 <nil>
}

func ExampleCache_SetMulti() {
	cacher.Flush()
	multi := map[string][]byte{
		"key1": []byte("value1"),
		"key2": []byte("value2"),
	}

	cacher.SetMulti(multi, 0)
	val1, _, _ := cacher.Get("key1")
	val2, _, _ := cacher.Get("key2")

	fmt.Println(string(val1))
	fmt.Println(string(val2))

	// Output:
	// value1
	// value2
}

func ExampleCache_CompareAndReplace() {
	cacher.Flush()
	cacher.Set("key1", []byte("hello world"), 0)
	_, token, _ := cacher.Get("key1")

	var err error
	var val []byte
	err = cacher.CompareAndReplace(token+"FALSE", "key1", []byte("replacement1"), 0)
	val, _, _ = cacher.Get("key1")
	fmt.Println(err, string(val))

	err = cacher.CompareAndReplace(token, "key1", []byte("replacement2"), 0)
	val, _, _ = cacher.Get("key1")
	fmt.Println(err, string(val))

	// Output:
	// Key `key1` does not exist. hello world
	// <nil> replacement2
}

func ExampleCache_Replace() {
	cacher.Flush()
	var err error

	err = cacher.Replace("key1", []byte("replacement"), 0)
	fmt.Println(err)

	cacher.Set("key1", []byte("value1"), 0)
	err = cacher.Replace("key1", []byte("replacement"), 0)
	fmt.Println(err)

	// Output:
	// Key `key1` does not exist.
	// <nil>
}

func ExampleCache_Get() {
	cacher.Flush()
	var value []byte
	var token string
	var err error

	value, token, err = cacher.Get("non-existing")
	fmt.Println(string(value), token, err)

	cacher.Set("key1", []byte("Hello world!"), 0)
	value, token, err = cacher.Get("key1")
	fmt.Println(string(value), token, err)

	// Output:
	// Key `non-existing` does not exist.
	// Hello world! 86fb269d190d2c85f6e0468ceca42a20 <nil>
}

func ExampleCache_GetMulti() {
	cacher.Flush()
	multi := map[string][]byte{
		"key1": []byte("value1"),
		"key2": []byte("value2"),
	}

	cacher.SetMulti(multi, 0)

	keys := []string{"key1", "key2"}

	values, tokens, bools := cacher.GetMulti(keys)
	fmt.Println(values["key1"], values["key2"])
	fmt.Println(tokens["key1"], tokens["key2"])
	fmt.Println(bools["key1"], bools["key2"])

	// Output:
	// [118 97 108 117 101 49] [118 97 108 117 101 50]
	// 9946687e5fa0dab5993ededddb398d2e f066ce9385512ee02afc6e14d627e9f2
	// <nil> <nil>
}

func ExampleCache_Increment() {
	cacher.Flush()
	cacher.Increment("key1", 0, 1, 0)
	cacher.Increment("key1", 0, 1, 0)
	v, _, _ := cacher.Get("key1")

	num, _ := encoding.BytesInt64(v)
	fmt.Println(num)

	cacher.Set("key2", []byte("string value, not incrementable"), 0)
	ok := cacher.Increment("key2", 0, 5, 0)
	v2, _, _ := cacher.Get("key2")
	fmt.Println(ok, string(v2))

	// Output:
	// 1
	// Value for key `key2` could not be encoded. string value, not incrementable
}

func ExampleCache_Decrement() {
	cacher.Flush()
	cacher.Decrement("key1", 10, 1, 0)
	cacher.Decrement("key1", 10, 3, 0)
	v, _, _ := cacher.Get("key1")

	num, _ := encoding.BytesInt64(v)
	fmt.Println(num)

	cacher.Set("key2", []byte("string value, not decrementable"), 0)
	ok := cacher.Decrement("key2", 0, 5, 0)
	v2, _, _ := cacher.Get("key2")
	fmt.Println(ok, string(v2))

	// Output:
	// 7
	// Value for key `key2` could not be encoded. string value, not decrementable
}

func ExampleCache_Flush() {
	cacher.Flush()
	var errs map[string]error

	keys := []string{"key1", "key2"}
	multi := map[string][]byte{
		"key1": []byte("value1"),
		"key2": []byte("value2"),
	}
	cacher.SetMulti(multi, 0)

	_, _, errs = cacher.GetMulti(keys)
	fmt.Println(errs["key1"], errs["key2"])

	cacher.Flush()
	_, _, errs = cacher.GetMulti(keys)
	fmt.Println(errs["key1"], errs["key2"])

	// Output:
	// <nil> <nil>
	// Key `key1` does not exist. Key `key2` does not exist.
}

func ExampleCache_Delete() {
	cacher.Flush()
	var err error

	cacher.Set("key1", []byte("value1"), 0)
	err = cacher.Delete("key1")
	fmt.Println(err)

	err = cacher.Delete("non-existing")
	fmt.Println(err)

	// Output:
	// <nil>
	// Key `non-existing` was not found.
}

func ExampleCache_DeleteMulti() {
	cacher.Flush()
	keys := []string{"key1", "key2", "non-existing"}

	multi := map[string][]byte{
		"key1": []byte("value1"),
		"key2": []byte("value2"),
	}
	cacher.SetMulti(multi, 0)

	oks := cacher.DeleteMulti(keys)
	fmt.Println(oks["key1"])
	fmt.Println(oks["key2"])
	fmt.Println(oks["non-existing"])

	// Output:
	// <nil>
	// <nil>
	// Key `non-existing` was not found.
}

func ExampleCache_Touch() {
	cacher.Flush()
	cacher.Set("key1", []byte("value1"), 1)
	err := cacher.Touch("key1", 5)
	fmt.Println(err)

	// Output:
	// <nil>
}
