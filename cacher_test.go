package cacher_test

import (
	"testing"

	"github.com/jelmersnoeck/cacher"
)

func acceptCacher(c cacher.Cacher) bool {
	return true
}

func TestMemoryCache(t *testing.T) {
	cache := cacher.NewMemoryCache(0)

	if !acceptCacher(cache) {
		t.Errorf("Expected MemoryCache to be accepted")
		t.FailNow()
	}
}
