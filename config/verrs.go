package config

import "strings"

type ValidationErrors []error

func (verrs ValidationErrors) HasErrors() bool {
	return len(verrs) > 0
}

func (verrs ValidationErrors) Error() string {
	errs := make([]string, len(verrs))

	for i := range verrs {
		errs[i] = verrs[i].Error()
	}
	return strings.Join(errs, ", ")
}
