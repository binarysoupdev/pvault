package data

import (
	"encoding/binary"
	"encoding/json"
	"os"
	"pvault/errors"

	"github.com/binarysoupdev/cryptool/crypt"
)

func SaveVersionedRecord(path string, version uint16, data []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return errors.Chain(err, "error creating record file")
	}
	defer file.Close()

	header := make([]byte, 2)
	binary.BigEndian.PutUint16(header, version)

	file.Write(header)
	file.Write(data)

	return nil
}

func SaveEncryptedRecord[T any](path string, password string, version uint16, header []byte, record T) error {
	c, salt := crypt.NewFromPassword(password)

	plaintext, err := json.Marshal(record)
	if err != nil {
		return errors.Chain(err, "error marshaling json")
	}

	ciphertext := c.Encrypt(plaintext)

	data := make([]byte, len(header)+len(salt)+len(ciphertext))
	copy(data, header)
	copy(data[len(header):], salt)
	copy(data[len(header)+len(salt):], ciphertext)

	return SaveVersionedRecord(path, version, data)
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
