# Cacher

[![TravisCI](https://travis-ci.org/jelmersnoeck/cacher.svg)](https://travis-ci.org/jelmersnoeck/cacher) [![GoDoc](https://godoc.org/github.com/jelmersnoeck/cacher?status.svg)](https://godoc.org/github.com/jelmersnoeck/cacher)

Cacher is a port of the PHP library [Scrapbook](https://github.com/matthiasmullie/scrapbook).

It defines an interface to interact with several cache systems without having to
worry about the implementation of said cache layer.

## Usage

See the specific packages' GoDoc reference for examples.

## Implementations

### Memory

[![GoDoc](https://godoc.org/github.com/jelmersnoeck/cacher/memory?status.svg)](https://godoc.org/github.com/jelmersnoeck/cacher/memory)

Memory stores all the data in memory. This is a non persistent cache store that
will be flushed every time the application that uses the cache is terminates.

This cache is perfect to use for testing. There are no other dependencies
required other than enough available memory.

### Redis

[![GoDoc](https://godoc.org/github.com/jelmersnoeck/cacher/redis?status.svg)](https://godoc.org/github.com/jelmersnoeck/cacher/redis)

Redis stores all the data in a Redis instance. This cache relies on the
`github.com/garyburd/redigo/redis` package to communicate with Redis.
