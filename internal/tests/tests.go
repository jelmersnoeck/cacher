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

func NotPresent(t *testing.T, cache cacher.Cacher, key string) {
	if _, _, ok := cache.Get(key); ok {
		FailMsg(t, cache, "Expected `"+key+"` not to be present")
	}
}

func FailMsg(t *testing.T, cache cacher.Cacher, msg string) {
	errMsg := reflect.TypeOf(cache).String() + ": " + msg
	t.Errorf(errMsg)
	t.FailNow()
}
