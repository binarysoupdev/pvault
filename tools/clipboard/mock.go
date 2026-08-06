package clipboard

type MockClipboard struct {
	Data             string
	UnsupportedError error
	ReadError        error
	WriteError       error
}

func Mock() *MockClipboard {
	return &MockClipboard{}
}

func (c MockClipboard) CheckUnsupported() error {
	return c.UnsupportedError
}

func (c MockClipboard) Read() (string, error) {
	return c.Data, c.ReadError
}

func (c *MockClipboard) Write(data string) error {
	c.Data = data
	return c.WriteError
}
