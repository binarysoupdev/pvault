package errors

import "fmt"

func Chain(err error, msg string) error {
	return fmt.Errorf("%s\n  %s", msg, err.Error())
}

func ChainFormat(err error, format string, a ...any) error {
	return Chain(err, fmt.Sprintf(format, a...))
}
