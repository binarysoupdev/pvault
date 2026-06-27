package vault

import "fmt"

type OutOfDateError struct {
	Version int
}

func newOutOfDateError(version int) OutOfDateError {
	return OutOfDateError{
		Version: version,
	}
}

func (err OutOfDateError) Error() string {
	return fmt.Sprintf("version \"%d\" out-of-date", err.Version)
}
