// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package memory_test

import (
	"testing"

	"github.com/jelmersnoeck/cacher/internal/tests"
	"github.com/jelmersnoeck/cacher/memory"
)

func TestLimit(t *testing.T) {
	cache := memory.New(30)

	cache.Add("key1", []byte("value1"), 0)
	cache.Add("key2", []byte("value2"), 0)
	cache.Add("key3", []byte("value3"), 0)
	cache.Add("key4", []byte("value4"), 0)
	cache.Add("key5", []byte("value5"), 0)

	tests.Compare(t, cache, "key1", "value1")
	tests.Compare(t, cache, "key2", "value2")
	tests.Compare(t, cache, "key3", "value3")
	tests.Compare(t, cache, "key4", "value4")
	tests.Compare(t, cache, "key5", "value5")

	cache.Add("key6", []byte("value6"), 0)

	tests.NotPresent(t, cache, "key1")

	cache.Delete("key3")
	cache.Add("key7", []byte("value7"), 0)
	tests.Compare(t, cache, "key2", "value2")

	cache.Add("key8", []byte("value8"), 0)

	// This is key5 due to the fact that we fetched key2 before. This pushed it
	// back as the most active key, so it wouldn't be deleted immediately.
	tests.NotPresent(t, cache, "key4")

	cache.Add("key9", []byte("value9"), 0)

	tests.NotPresent(t, cache, "key5")
}
