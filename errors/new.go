package errors

import (
	"fmt"
)

func New(msg string) error {
	return fmt.Errorf("%s", msg)
}

func Format(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}
