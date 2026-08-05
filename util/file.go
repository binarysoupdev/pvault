package util

import (
	"io"
	"os"

	"github.com/binarysoupdev/go-commando/errors"
)

func CreateEmptyFile(path string) error {
	return os.WriteFile(path, []byte{}, 0666)
}

func CopyFile(dest string, src string) error {
	s, err := os.Open(src)
	if err != nil {
		return errors.Chain(err, "error opening source file")
	}
	defer s.Close()

	d, err := os.Create(dest)
	if err != nil {
		return errors.Chain(err, "error creating destination file")
	}
	defer d.Close()

	_, err = io.Copy(d, s)
	if err != nil {
		return errors.Chain(err, "error copying data")
	}

	return nil
}
