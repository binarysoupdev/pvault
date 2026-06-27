package index

const CURRENT_VERSION = 1

type Codec struct{}

func (c Codec) Version() int {
	return CURRENT_VERSION
}
