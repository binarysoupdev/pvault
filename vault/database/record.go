package database

import (
	"encoding/json"
	"os"

	"github.com/binarysoupdev/go-commando/errors"

	"github.com/binarysoupdev/cryptool/crypt"
)

func SaveEncryptedRecord[T any](path string, password string, header []byte, record T) error {
	c, salt := crypt.NewFromPassword(password)

	plaintext, err := json.Marshal(record)
	if err != nil {
		return errors.Chain(err, "error marshaling json")
	}

	ciphertext := c.Encrypt(plaintext)

	file, err := os.Create(path)
	if err != nil {
		return errors.Chain(err, "error creating record file")
	}
	defer file.Close()

	file.Write(header)
	file.Write(salt)
	file.Write(ciphertext)

	return nil
}

func DecryptRecord[T any](password string, ciphertext []byte) (T, error) {
	var obj T
	c := crypt.LoadFromPassword(password, ciphertext[:crypt.SALT_SIZE])

	plaintext, err := c.Decrypt(ciphertext[crypt.SALT_SIZE:])
	if err != nil {
		return obj, errors.Chain(err, "error decrypting ciphertext")
	}

	err = json.Unmarshal(plaintext, &obj)
	if err != nil {
		return obj, errors.Chain(err, "error unmarshaling json")
	}

	return obj, nil
}
