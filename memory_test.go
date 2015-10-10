package cacher_test

import (
	"testing"

	"github.com/jelmersnoeck/cacher"
)

func TestMemorySet(t *testing.T) {
	cache := cacher.NewMemoryCache(0)

	if !cache.Set("key1", "value", 0) {
		t.Errorf("Expecting `key1` to be `value`")
		t.Fail()
	}

	if !cache.Set("key2", 2, 0) {
		t.Errorf("Expecting `key2` to be `2`")
		t.Fail()
	}
}

func TestMemoryAdd(t *testing.T) {
	cache := cacher.NewMemoryCache(0)

	if !cache.Add("key1", "value1", 0) {
		t.Errorf("Expecting `key1` to be added to the cache")
		t.FailNow()
	}

	if cache.Add("key1", "value2", 0) {
		t.Errorf("Expecting `key1` not to be added to the cache")
		t.FailNow()
	}

	compare(t, cache, "key1", "value1")
}

func TestMemoryReplace(t *testing.T) {
	cache := cacher.NewMemoryCache(0)

	if cache.Replace("key1", "value1", 0) {
		t.Errorf("Key1 is not set yet, should not be able to replace.")
		t.FailNow()
	}

	cache.Set("key1", "value1", 0)
	if !cache.Replace("key1", "value1", 0) {
		t.Errorf("Key1 has been set, should be able to replace.")
		t.FailNow()
	}
}

func TestMemoryGet(t *testing.T) {
	cache := cacher.NewMemoryCache(0)

	cache.Set("key1", "value1", 0)
	compare(t, cache, "key1", "value1")

	if _, ok := cache.Get("key2"); ok {
		t.Errorf("Key2 is not present, ok should be false.")
		t.FailNow()
	}
}

func TestMemoryFlush(t *testing.T) {
	cache := cacher.NewMemoryCache(0)

	cache.Set("key1", "value1", 0)
	compare(t, cache, "key1", "value1")

	if !cache.Flush() {
		t.Fail()
	}

	if v, _ := cache.Get("key1"); v != nil {
		t.Errorf("Expecting `key1` to be nil")
		t.Fail()
	}
}

func TestLimit(t *testing.T) {
	cache := cacher.NewMemoryCache(250)

	cache.Add("key1", "value1", 0)
	cache.Add("key2", "value2", 0)
	cache.Add("key3", "value3", 0)
	cache.Add("key4", "value4", 0)
	cache.Add("key5", "value5", 0)

	compare(t, cache, "key1", "value1")
	compare(t, cache, "key2", "value2")
	compare(t, cache, "key3", "value3")
	compare(t, cache, "key4", "value4")
	compare(t, cache, "key5", "value5")

	cache.Add("key6", "value6", 0)

	notPresent(t, cache, "key1")

	cache.Delete("key3")
	cache.Add("key7", "value7", 0)
	compare(t, cache, "key2", "value2")

	cache.Add("key8", "value8", 0)

	// This is key5 due to the fact that we fetched key2 before. This pushed it
	// back as the most active key, so it wouldn't be deleted immediately.
	notPresent(t, cache, "key4")

	cache.Add("key9", "value9", 0)

	notPresent(t, cache, "key5")
}

func compare(t *testing.T, cache cacher.Cacher, key string, value interface{}) {
	if v, _ := cache.Get(key); v != value {
		t.Errorf("Expected `" + key + "` to equal `" + value.(string) + "`")
		t.FailNow()
	}
}

func notPresent(t *testing.T, cache cacher.Cacher, key string) {
	if _, ok := cache.Get(key); ok {
		t.Errorf("Expected `" + key + "` not to be present")
		t.FailNow()
	}
}
