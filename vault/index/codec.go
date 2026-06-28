package index

type Codec struct{}

func (Codec) Version() uint16 {
	return 2
}
