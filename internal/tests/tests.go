// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

// Package tests provides helper methods for testing
package tests

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/jelmersnoeck/cacher"
	"github.com/jelmersnoeck/cacher/internal/encoding"
)

// Compare compares a cached value given by the cache and the key to the value
// passed in as an interface. If the values do not match, the given test will
// receive a `FailNow()` call and print an appropiate error message.
func Compare(t *testing.T, cache cacher.Cacher, key string, value interface{}) {
	_, ok := value.(int)
	if ok {
		val := int64(value.(int))
		v, _, _ := cache.Get(key)
		valInt, _ := encoding.BytesInt64(v)
		if valInt != val {
			msg := "Expected `" + key + "` to equal `" + strconv.FormatInt(val, 10) + "`, is `" + strconv.FormatInt(valInt, 10) + "`"
			FailMsg(t, cache, msg)
		}
	} else {
		value = value.(string)
		if v, _, _ := cache.Get(key); string(v) != value {
			msg := "Expected `" + key + "` to equal `" + value.(string) + "`"
			FailMsg(t, cache, msg)
		}
	}
}

// NotPresent ensures that the given key is not present in the given cache. If
// it is present, the test will fail and print an error message.
func NotPresent(t *testing.T, cache cacher.Cacher, key string) {
	if _, _, err := cache.Get(key); err != nil {
		FailMsg(t, cache, "Expected `"+key+"` not to be present")
	}
}

// FailMsg will print an error message that specifies the cache type and fail
// the given test instance.
func FailMsg(t *testing.T, cache cacher.Cacher, msg string, s ...interface{}) {
	errMsg := reflect.TypeOf(cache).String() + ": " + msg
	t.Errorf(errMsg, s...)
	t.FailNow()
}
