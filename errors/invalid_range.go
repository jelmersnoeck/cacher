package errors

import "fmt"

// InvalidRange errors are used when we want to increment or decrement a value
// which is out of range of the specific operation.
type InvalidRange struct {
	initial int64
	offset  int64
}

func (e InvalidRange) Error() string {
	return fmt.Sprintf(
		"The range `%s` to `%s` is not supported.",
		e.initial,
		e.offset,
	)
}

func NewInvalidRange(initial, offset int64) error {
	return InvalidRange{
		initial: initial,
		offset:  offset,
	}
}
