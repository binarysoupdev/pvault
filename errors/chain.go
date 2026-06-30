package errors

import "fmt"

func Chain(err error, msg string) error {
	return fmt.Errorf("%s\n  %s", msg, err.Error())
}

//TODO: add ChainFormat
