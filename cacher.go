// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

// Package cacher provides a uniform interface for different caching strategies.
package cacher

import "time"

// Cacher is the Caching interface that uniforms all the different strategies.
type Cacher interface {
	Add(key string, value interface{}, ttl int) bool
	Set(key string, value interface{}, ttl int) bool
	Delete(key string) bool
	Get(key string) interface{}
	Flush() bool
}

type TimeNowF func() time.Time

var TimeNow TimeNowF

func init() {
	ResetTimeNow()
}

func ResetTimeNow() {
	TimeNow = func() time.Time {
		return time.Now()
	}
}
