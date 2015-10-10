package cacher_test

import (
	"testing"

	"github.com/jelmersnoeck/cacher"
)

func TestMemorySet(t *testing.T) {
	cache := cacher.NewMemoryCache()

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
	cache := cacher.NewMemoryCache()

	if !cache.Add("key1", "value1", 0) {
		t.Errorf("Expecting `key1` to be added to the cache")
		t.FailNow()
	}

	if cache.Add("key1", "value2", 0) {
		t.Errorf("Expecting `key1` not to be added to the cache")
		t.FailNow()
	}

	if cache.Get("key1") != "value1" {
		t.Errorf("Expecting `key1` to equal `value1`")
		t.FailNow()
	}
}

func TestMemoryGet(t *testing.T) {
	cache := cacher.NewMemoryCache()

	cache.Set("key1", "value", 0)
	if cache.Get("key1") != "value" {
		t.Errorf("Expecting `key1` to be `value`")
		t.Fail()
	}
}

func TestMemoryFlush(t *testing.T) {
	cache := cacher.NewMemoryCache()

	cache.Set("key1", "value1", 0)
	if cache.Get("key1") != "value1" {
		t.Errorf("Expecting `key1` to equal `value1`")
		t.Fail()
	}

	if !cache.Flush() {
		t.Fail()
	}

	if cache.Get("key1") != nil {
		t.Errorf("Expecting `key1` to be nil")
		t.Fail()
	}
}
