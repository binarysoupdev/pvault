package data

type NotSupportedError struct{}

func (NotSupportedError) Error() string {
	return "not supported"
}
