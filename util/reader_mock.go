package util

type MockReader struct {
	callCount  int
	ReadErrors []error
}

func (m *MockReader) Read(_ []byte) (int, error) {
	m.callCount += 1
	return 0, m.ReadErrors[min(m.callCount, len(m.ReadErrors))-1]
}
