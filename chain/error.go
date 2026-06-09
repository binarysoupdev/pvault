package chain

import "fmt"

func New(msg string) error {
	return fmt.Errorf("%s", msg)
}

func Error(err error, msg string) error {
	return fmt.Errorf("%s\n  %s", msg, err.Error())
}
