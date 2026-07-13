package common

type NotSupportedError struct{}

func (NotSupportedError) Error() string {
	return "not supported"
}
