package chain

import "fmt"

func Error(err error, msg string) error {
	return fmt.Errorf("%s\n  %s", msg, err.Error())
}
