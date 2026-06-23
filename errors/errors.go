package errors

type Errors []error

func (errs *Errors) Append(err error) {
	*errs = append(*errs, err)
}

func (errs Errors) Collapse() error {
	if len(errs) == 0 {
		return nil
	}

	const SEPARATOR = ", "
	msg := ""

	for _, err := range errs {
		msg += err.Error() + SEPARATOR
	}

	return New(msg[:len(msg)-len(SEPARATOR)])
}
