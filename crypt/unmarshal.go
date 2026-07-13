package crypt

import (
	"encoding/json"

	"github.com/binarysoupdev/cryptool/crypt"
	"github.com/binarysoupdev/go-commando/errors"
)

func Unmarshal[T any](password string, ciphertext Ciphertext) (T, error) {
	var obj T
	c := crypt.LoadFromPassword(password, ciphertext.Salt())

	plaintext, err := c.Decrypt(ciphertext.Text())
	if err != nil {
		return obj, errors.Chain(err, "error decrypting ciphertext")
	}

	err = json.Unmarshal(plaintext, &obj)
	if err != nil {
		return obj, errors.Chain(err, "error unmarshaling json")
	}

	return obj, nil
}
