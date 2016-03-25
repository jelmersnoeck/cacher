package errors

import "fmt"

// AlreadyExisingKey is an error type used for when the key is expected not to
// exist in the cache yet, but it is already present.
type AlreadyExistingKey struct {
	key string
}

func (e AlreadyExistingKey) Error() string {
	return fmt.Sprintf("Key `%s` already exists.", e.key)
}

func NewAlreadyExistingKey(key string) error {
	return AlreadyExistingKey{
		key: key,
	}
}
