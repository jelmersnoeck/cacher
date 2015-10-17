# Cacher

[![TravisCI](https://travis-ci.org/jelmersnoeck/cacher.svg)](https://travis-ci.org/jelmersnoeck/cacher) [![GoDoc](https://godoc.org/github.com/jelmersnoeck/cacher?status.svg)](https://godoc.org/github.com/jelmersnoeck/cacher)

Cacher is a port of the PHP library [Scrapbook](https://github.com/matthiasmullie/scrapbook).

It defines an interface to interact with several cache systems without having to
worry about the implementation of said cache layer.

## Implementations

### Memory

[![GoDoc](https://godoc.org/github.com/jelmersnoeck/cacher/memory?status.svg)](https://godoc.org/github.com/jelmersnoeck/cacher/memory)

Memory stores all the data in memory. This is a non persistent cache store that
will be flushed every time the application that uses the cache is terminates.

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

### Redis

[![GoDoc](https://godoc.org/github.com/jelmersnoeck/cacher/redis?status.svg)](https://godoc.org/github.com/jelmersnoeck/cacher/redis)

Redis stores all the data in a Redis instance. This cache relies on the
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
