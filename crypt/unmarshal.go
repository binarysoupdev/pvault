package crypt

import (
	"encoding/json"
	"io"

	"github.com/binarysoupdev/cryptool/crypt"
)

func Unmarshal[T any](password string, ciphertext Ciphertext) (T, error) {
	var obj T
	c := crypt.LoadFromPassword(password, ciphertext.Salt())

	plaintext, err := c.Decrypt(ciphertext.Text())
	if err != nil {
		return obj, err
	}

	err = json.Unmarshal(plaintext, &obj)
	if err != nil {
		return obj, err
	}

	return obj, nil
}

func Decode[T any](r io.Reader, password string) (T, error) {
	var obj T
	ciphertext, err := io.ReadAll(r)
	if err != nil {
		return obj, err
	}

	return Unmarshal[T](password, ciphertext)
}
