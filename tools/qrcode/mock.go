package qrcode

type MockRenderer struct {
	Text        string
	RenderError error
}

func Mock() *MockRenderer {
	return &MockRenderer{}
}

func (m *MockRenderer) RenderToStdout(text string) error {
	m.Text = text
	return m.RenderError
}
