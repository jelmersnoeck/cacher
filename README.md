# Cacher

Cacher is a port of the PHP a library [Scrapbook](https://github.com/matthiasmullie/scrapbook).

It defines an interface to interact with several cache systems without having to
worry about the implementation of said cache layer.

## Impelmentations

### MemoryCache

MemoryCache stores all the data in memory. This is a non persistent cache store
that will be flushed every time the application that uses the cache is
terminates.

This cache is perfect to use for testing. There are no other dependencies
required other than enough available memory.
