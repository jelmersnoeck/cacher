package errors

import "fmt"

// Encoding is an error used when encoding data to bytes fails.
type Encoding struct {
	key string
}

func (e Encoding) Error() string {
	return fmt.Sprintf("Value for key `%s` could not be encoded.", e.key)
}

func NewEncoding(key string) error {
	return Encoding{
		key: key,
	}
}
