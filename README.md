# Cacher

[![TravisCI](https://travis-ci.org/jelmersnoeck/cacher.svg)](https://travis-ci.org/jelmersnoeck/cacher) [![GoDoc reference](https://camo.githubusercontent.com/fb9e66520f8775e97dcacdf366d0dee7828df53f/68747470733a2f2f676f646f632e6f72672f6769746875622e636f6d2f676f2d6d617274696e692f6d617274696e693f7374617475732e706e67)](https://godoc.org/github.com/jelmersnoeck/cacher)

Cacher is a port of the PHP a library [Scrapbook](https://github.com/matthiasmullie/scrapbook).

It defines an interface to interact with several cache systems without having to
worry about the implementation of said cache layer.

## Methods

### `Add(key string, value interface{}, ttl int) bool`

Adds a new key to the cache if the key is not already stored. If the key is
already stored false will be returned.

### `Set(key string, value interface{}, ttl int) bool`

Sets the value for the specified key, regardless of wether or not the key has
already been set. If the key has already been set, it will overwrite the
previous value.

### `Get(key string) interface{}`

Gets the value for the given key.

### `Flush() bool`

Resets the cache store and deletes all cached values.

### `Delete(key string) bool`

Deletes a key from the cache.

## Implementations

### MemoryCache

MemoryCache stores all the data in memory. This is a non persistent cache store
that will be flushed every time the application that uses the cache is
terminates.

This cache is perfect to use for testing. There are no other dependencies
required other than enough available memory.
