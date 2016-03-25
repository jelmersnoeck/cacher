package errors

import "fmt"

// ValueBelowZero errors are used for incrementing issues where the value drops
// below zero, which is not supported by all caches.
type ValueBelowZero struct {
	key string
}

func (e ValueBelowZero) Error() string {
	return fmt.Sprintf("Value for key `%s` dropped below 0.", e.key)
}

func NewValueBelowZero(key string) error {
	return ValueBelowZero{
		key: key,
	}
}
