package crypt

import (
	"encoding/json"

	"github.com/binarysoupdev/cryptool/crypt"
	"github.com/binarysoupdev/go-commando/errors"
)

func Marshal[T any](password string, obj T) (Ciphertext, error) {
	plaintext, err := json.Marshal(obj)
	if err != nil {
		return nil, errors.Chain(err, "error marshaling json")
	}

	c, salt := crypt.NewFromPassword(password)
	ciphertext := c.Encrypt(plaintext)

	return append(salt, ciphertext...), nil
}
