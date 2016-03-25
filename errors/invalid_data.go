package errors

import "fmt"

// InvalidData errors are ussed when we have data for a given key, but the data
// type does not match a slice of bytes (`[]byte`).
type InvalidData struct {
	key string
}

func (e InvalidData) Error() string {
	return fmt.Sprintf("Key `%s` was not of valid data type.", e.key)
}

func NewInvalidData(key string) error {
	return InvalidData{
		key: key,
	}
}
