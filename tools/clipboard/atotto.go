package clipboard

import "github.com/atotto/clipboard"

type AtottoClipboard struct{}

func (AtottoClipboard) CheckUnsupported() error {
	if clipboard.Unsupported {
		return clipboard.WriteAll("")
	}
	return nil
}

func (AtottoClipboard) Read() (string, error) {
	return clipboard.ReadAll()
}

func (AtottoClipboard) Write(data string) error {
	return clipboard.WriteAll(data)
}
