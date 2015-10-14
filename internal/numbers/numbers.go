// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

// Package numbers provides helper methods to deal with numbers.
package numbers

import (
	"bytes"
	"strconv"
)

// BytesString converts a slice of bytes into a string. It takes care of dealing
// with 0 bytes by checking when the first occurance happens and cutting of the
// slice at that point.
func BytesString(bts []byte) string {
	n := bytes.IndexByte(bts, 0)

	if n == -1 {
		n = len(bts)
	}

	return string(bts[:n])
}

// Int64Bytes converts an int64 value into a string and returns that as a slice
// of bytes.
func Int64Bytes(val int64) []byte {
	intStr := strconv.FormatInt(int64(val), 10)

	return []byte(intStr)
}

// BytesInt64 converts a previously converted int64 to string byte slice into
// its int64 equivalent. If the byte slice doesn't turn out to be a correct
// int64 value, it will return false.
func BytesInt64(bts []byte) (int64, bool) {
	btsStr := BytesString(bts)

	var ok bool
	val, err := strconv.ParseInt(btsStr, 10, 64)

	if err == nil {
		ok = true
	}

	return int64(val), ok
}
