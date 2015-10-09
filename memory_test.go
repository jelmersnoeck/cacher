package cacher_test

import (
	"testing"

	"github.com/jelmersnoeck/cacher"
)

func TestMemorySet(t *testing.T) {
	cache := cacher.NewMemoryCache()

	if !cache.Set("key1", "value") {
		t.Fail()
	}

	if !cache.Set("key2", 2) {
		t.Fail()
	}
}

func TestMemoryGet(t *testing.T) {
	cache := cacher.NewMemoryCache()

	cache.Set("key1", "value")
	if cache.Get("key1") != "value" {
		t.Fail()
	}
}
