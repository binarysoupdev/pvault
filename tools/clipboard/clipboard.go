package clipboard

type Clipboard interface {
	CheckUnsupported() error
	Read() (string, error)
	Write(data string) error
}
