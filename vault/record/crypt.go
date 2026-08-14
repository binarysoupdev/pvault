package record

import (
	"encoding/json"

	"github.com/binarysoupdev/cryptool/crypt"
	"github.com/binarysoupdev/go-extensions/errors"
)

func Encrypt(r Record, password string) ([]byte, error) {
	plaintext, err := json.Marshal(r)
	if err != nil {
		return nil, errors.Chain(err, "error marshaling json")
	}

	if password == "" {
		return plaintext, nil
	}

	c, salt := crypt.NewFromPassword(password)
	ciphertext := c.Encrypt(plaintext)

	return append(salt, ciphertext...), nil
}

func Decrypt[T Record](bytes []byte, password string) (T, error) {
	var r T

	plaintext, err := decryptBytes(bytes, password)
	if err != nil {
		return r, err
	}

	err = json.Unmarshal(plaintext, &r)
	if err != nil {
		return r, errors.Chain(err, "error unmarshaling json")
	}

	return r, nil
}

func decryptBytes(ciphertext []byte, password string) ([]byte, error) {
	if password == "" {
		return ciphertext, nil
	}

	if len(ciphertext) < crypt.SALT_SIZE+1 {
		return nil, errors.Format("length too short: %d", len(ciphertext))
	}

	c := crypt.LoadFromPassword(password, ciphertext[:crypt.SALT_SIZE])
	return c.Decrypt(ciphertext[crypt.SALT_SIZE:])
}
