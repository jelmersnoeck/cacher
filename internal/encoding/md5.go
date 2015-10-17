// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

// Package encoding provides helper methods to deal with encoding data.
package encoding

import (
	"crypto/md5"
	"encoding/hex"
)

// Md5Sum converts an array of bytes to an md5 Sum string.
func Md5Sum(value []byte) string {
	hasher := md5.New()
	hasher.Write(value)
	return hex.EncodeToString(hasher.Sum(nil))
}
