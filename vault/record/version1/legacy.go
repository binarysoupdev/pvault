package v1

import (
	"io"
	"os"
	"pvault/crypt"
)

const LEGACY_HASH_SIZE = 60

func EncodeFromLegacy(w io.Writer, name string, bytes []byte) {
	writeHeader(w, name)
	w.Write(bytes[LEGACY_HASH_SIZE:])
}

func (r Record) MarshalToLegacy(path string, password string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	hash := make([]byte, LEGACY_HASH_SIZE)
	file.Write(hash)

	_, err = crypt.Encode(file, password, r)
	return err
}
