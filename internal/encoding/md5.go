// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

// Package encoding provides helper methods to deal with encoding data.
package encoding

import "crypto/md5"

// Md5Sum converts an array of bytes to an md5 Sum string.
func Md5Sum(value []byte) string {
	md5Sum := md5.Sum(value)
	return string(md5Sum[:])
}
