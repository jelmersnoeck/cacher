package cacher_test

import (
	"testing"

	"github.com/jelmersnoeck/cacher"
	"github.com/jelmersnoeck/cacher/internal/tester"
)

func TestMemoryCollection(t *testing.T) {
	tester.RunCacher(t, cacher.NewMemoryCache(0))
}

func TestLimit(t *testing.T) {
	cache := cacher.NewMemoryCache(30)

	cache.Add("key1", []byte("value1"), 0)
	cache.Add("key2", []byte("value2"), 0)
	cache.Add("key3", []byte("value3"), 0)
	cache.Add("key4", []byte("value4"), 0)
	cache.Add("key5", []byte("value5"), 0)

	tester.Compare(t, cache, "key1", "value1")
	tester.Compare(t, cache, "key2", "value2")
	tester.Compare(t, cache, "key3", "value3")
	tester.Compare(t, cache, "key4", "value4")
	tester.Compare(t, cache, "key5", "value5")

	cache.Add("key6", []byte("value6"), 0)

	tester.NotPresent(t, cache, "key1")

	cache.Delete("key3")
	cache.Add("key7", []byte("value7"), 0)
	tester.Compare(t, cache, "key2", "value2")

	cache.Add("key8", []byte("value8"), 0)

	// This is key5 due to the fact that we fetched key2 before. This pushed it
	// back as the most active key, so it wouldn't be deleted immediately.
	tester.NotPresent(t, cache, "key4")

	cache.Add("key9", []byte("value9"), 0)

	tester.NotPresent(t, cache, "key5")
}
