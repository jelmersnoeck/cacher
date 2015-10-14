package cacher_test

import (
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/jelmersnoeck/cacher"
	"github.com/jelmersnoeck/cacher/internal/tester"
)

func TestRedisCollection(t *testing.T) {
	c, _ := redis.Dial("tcp", ":6379")
	defer c.Close()

	cache := cacher.NewRedisCache(c)
	tester.RunCacher(t, cache)
}
