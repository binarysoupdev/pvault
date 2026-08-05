package util

import (
	"io"

	"github.com/binarysoupdev/go-commando/errors"
)

func ReadBytes(r io.Reader, n int) ([]byte, error) {
	bytes := make([]byte, n)
	count, err := r.Read(bytes)

	if count < n {
		return bytes, errors.Format("too few bytes (expected: %d, actual: %d)", n, count)
	}
	return bytes, err
}

func WriteBytes(w io.Writer, b ...[]byte) error {
	for _, bytes := range b {
		if _, err := w.Write(bytes); err != nil {
			return err
		}
	}
	return nil
}
