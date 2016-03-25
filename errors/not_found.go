package errors

import "fmt"

// NotFound errors are used when we expect the data to be in the cache already
// but it isn't.
type NotFound struct {
	key string
}

func (e NotFound) Error() string {
	return fmt.Sprintf("Key `%s` was not found.", e.key)
}

func NewNotFound(key string) error {
	return NotFound{
		key: key,
	}
}
