# Cacher

[![TravisCI](https://travis-ci.org/jelmersnoeck/cacher.svg)](https://travis-ci.org/jelmersnoeck/cacher) [![GoDoc reference](https://camo.githubusercontent.com/fb9e66520f8775e97dcacdf366d0dee7828df53f/68747470733a2f2f676f646f632e6f72672f6769746875622e636f6d2f676f2d6d617274696e692f6d617274696e693f7374617475732e706e67)](https://godoc.org/github.com/jelmersnoeck/cacher)

Cacher is a port of the PHP library [Scrapbook](https://github.com/matthiasmullie/scrapbook).

It defines an interface to interact with several cache systems without having to
worry about the implementation of said cache layer.

## Methods

#### `Add(key string, value []byte, ttl int) bool`

Adds a new key to the cache if the key is not already stored. If the key is
already stored false will be returned.

#### `Set(key string, value []byte, ttl int) bool`

Sets the value for the specified key, regardless of wether or not the key has
already been set. If the key has already been set, it will overwrite the
previous value.

#### `SetMulti(items map[string][]byte, ttl int) map[string]bool`

A shorthand to set multiple key/value combinations at a time. This uses `Set`
internally to add the items to the cache.

#### `Increment(key string, initial, offset, ttl int) bool`

Increments the initial value - or cached value if present - by offset.

#### `Decrement(key string, initial, offset, ttl int) bool`

Decrements the initial value - or cached value if present - by offset.

#### `Replace(key string, value []byte, ttl int) bool`

Replace will update a value, only if it is present. If it is not present, false
will be returned.

#### `Get(key string) ([]byte, string, bool)`

Gets the value for the given key combined with a CompareAndReplace token. If the
value is not present, `false` will be returned.

#### `GetMulti(keys []string) (map[string][]byte, map[string]string, map[string]bool)`

Gets a list of values for a range of given keys. As with `Get()`, it will return
a map with keys and CompareAndReplace tokens. If an item doesn't exist, it will
return false.

#### `CompareAndReplace(token, key string, value []byte, ttl int64) bool`

Sees if the key exists in the cache, if it doesn't it will return false. If it
does, it will compare the token with the token in the store. If the tokens do
not match, false will be returned. If they do match, the value will be replaced
with a new value and true will be returned.

#### `Flush() bool`

Resets the cache store and deletes all cached values.

#### `Delete(key string) bool`

Deletes a key from the cache.

#### `DeleteMulti(keys []string) map[string]bool`

Deletes multiple keys at a time and returns with a result set to see if the
deletes were successful.

## Implementations

### MemoryCache

MemoryCache stores all the data in memory. This is a non persistent cache store
that will be flushed every time the application that uses the cache is
terminates.

This cache is perfect to use for testing. There are no other dependencies
required other than enough available memory.

#### Usage

```go
package main

import (
    "fmt"

    "github.com/jelmersnoeck/cacher/memory"
)

func main() {
	cache := memory.New(30)
	cache.Add("key1", []byte("value1"), 0)

    v, token, ok := cache.Get("key1")

    if ok {
        fmt.Println(v, token)
    } else {
        fmt.Println("Something went wrong")
    }
}
```

### RedisCache

RedisCache stores all the data in a Redis instance. This cache relies on the
`github.com/garyburd/redigo/redis` package to communicate with Redis.

#### Usage

```go
package main

import (
    "fmt"

	"github.com/garyburd/redigo/redis"
	rcache "github.com/jelmersnoeck/cacher/redis"
)

func main() {
	c, _ := redis.Dial("tcp", ":6379")
	cache := rcache.New(c)
	cache.Add("key1", []byte("value1"), 0)

    v, token, ok := cache.Get("key1")

    if ok {
        fmt.Println(v, token)
    } else {
        fmt.Println("Something went wrong")
    }
}
```
