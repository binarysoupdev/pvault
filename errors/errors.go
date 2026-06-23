package errors

type Errors []error

func (errs *Errors) Add(msg string) {
	errs.AddError(New(msg))
}

func (errs *Errors) AddError(err error) {
	*errs = append(*errs, err)
}

func (errs Errors) Collapse(sep string) error {
	if len(errs) == 0 {
		return nil
	}
	msg := ""

	for _, err := range errs {
		msg += err.Error() + sep
	}

	return New(msg[:len(msg)-len(sep)])
}
