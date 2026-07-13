package v1

import (
	"io"
	"pvault/crypt"
)

const LEGACY_HASH_SIZE = 60

func EncodeFromLegacy(w io.Writer, name string, bytes []byte) {
	writeHeader(w, name)
	w.Write(bytes[LEGACY_HASH_SIZE:])
}

func (r Record) EncodeToLegacy(w io.Writer, password string) error {
	hash := make([]byte, LEGACY_HASH_SIZE)
	w.Write(hash)

	_, err := crypt.Encode(w, password, r)
	return err
}
