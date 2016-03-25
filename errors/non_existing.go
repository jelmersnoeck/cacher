package errors

import "fmt"

// NonExistingKey errors are used when we look if a key exists. It is not to be
// confused with the `NotFound` error.
type NonExistingKey struct {
	key string
}

func (e NonExistingKey) Error() string {
	return fmt.Sprintf("Key `%s` does not exist.", e.key)
}

func NewNonExistingKey(key string) error {
	return NonExistingKey{
		key: key,
	}
}
