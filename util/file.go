package util

import "os"

func CreateEmptyFile(path string) error {
	return os.WriteFile(path, []byte{}, 0666)
}
