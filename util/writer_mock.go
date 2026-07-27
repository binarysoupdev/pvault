package util

type MockWriter struct {
	callCount   int
	WriteErrors []error
}

func (m *MockWriter) Write(_ []byte) (int, error) {
	m.callCount += 1
	return 0, m.WriteErrors[min(m.callCount, len(m.WriteErrors))-1]
}
