package errors

import (
	"errors"
	"fmt"
)

func New(msg string) error {
	return errors.New(msg)
}

func Format(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}
