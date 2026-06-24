package clipboard

type MockClipboard struct {
	Data        string
	Unsupported error
	Error       error
}

func Mock() *MockClipboard {
	return &MockClipboard{}
}

func (c MockClipboard) CheckUnsupported() error {
	return c.Unsupported
}

func (c MockClipboard) Read() (string, error) {
	return c.Data, c.Error
}

func (c *MockClipboard) Write(data string) error {
	c.Data = data
	return c.Error
}
