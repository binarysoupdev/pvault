package crypt

import (
	"encoding/json"
	"io"

	"github.com/binarysoupdev/cryptool/crypt"
)

func Marshal[T any](password string, obj T) (Ciphertext, error) {
	plaintext, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	c, salt := crypt.NewFromPassword(password)
	ciphertext := c.Encrypt(plaintext)

	return append(salt, ciphertext...), nil
}

func Encode[T any](w io.Writer, password string, obj T) (int, error) {
	ciphertext, err := Marshal(password, obj)
	if err != nil {
		return 0, err
	}

	return w.Write(ciphertext)
}
