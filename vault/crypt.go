package vault

import (
	"encoding/json"
	"os"
	"pvault/errors"
	"pvault/vault/record"

	"github.com/binarysoupdev/cryptool/crypt"
	"github.com/google/uuid"
)

func (v Vault) saveEncryptedRecord(r record.Record, password string) error {
	c, salt := crypt.NewFromPassword(password)

	plaintext, err := json.Marshal(r)
	if err != nil {
		return errors.Chain(err, "error marshaling json")
	}

	ciphertext := c.Encrypt(plaintext)

	err = os.WriteFile(v.RecordPath(r.ID), append(salt, ciphertext...), 0666)
	if err != nil {
		return errors.Chain(err, "error writing record file")
	}

	return nil
}

func (v Vault) loadEncryptedRecord(id uuid.UUID, password string) (record.Record, error) {
	r := record.Record{}

	raw, err := os.ReadFile(v.RecordPath(id))
	if err != nil {
		return r, errors.Chain(err, "error reading record file")
	}

	c := crypt.LoadFromPassword(password, raw[:crypt.SALT_SIZE])

	plaintext, err := c.Decrypt(raw[crypt.SALT_SIZE:])
	if err != nil {
		return r, errors.Chain(err, "error decrypting ciphertext")
	}

	err = json.Unmarshal(plaintext, &r)
	if err != nil {
		return r, errors.Chain(err, "error unmarshaling json")
	}

	return r, nil
}
