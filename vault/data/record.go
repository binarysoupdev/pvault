package data

import (
	"encoding/json"
	"os"
	"path/filepath"
	"pvault/errors"
	"pvault/vault/record"

	"github.com/binarysoupdev/cryptool/crypt"
	"github.com/google/uuid"
)

func (db DatabaseV2) RecordPath(id uuid.UUID) string {
	return filepath.Join(filepath.Dir(db.Path), id.String())
}

func (db DatabaseV2) SaveRecord(r record.Record, password string) error {
	c, salt := crypt.NewFromPassword(password)

	plaintext, err := json.Marshal(r)
	if err != nil {
		return errors.Chain(err, "error marshaling json")
	}

	ciphertext := c.Encrypt(plaintext)

	err = os.WriteFile(db.RecordPath(r.ID), append(salt, ciphertext...), 0666)
	if err != nil {
		return errors.Chain(err, "error writing record file")
	}

	return nil
}

func (db DatabaseV2) LoadRecord(id uuid.UUID, password string) (record.Record, error) {
	raw, err := os.ReadFile(db.RecordPath(id))
	if err != nil {
		return record.Record{}, errors.Chain(err, "error reading record file")
	}

	c := crypt.LoadFromPassword(password, raw[:crypt.SALT_SIZE])

	plaintext, err := c.Decrypt(raw[crypt.SALT_SIZE:])
	if err != nil {
		return record.Record{}, errors.Chain(err, "error decrypting ciphertext")
	}

	r := record.Record{}

	err = json.Unmarshal(plaintext, &r)
	if err != nil {
		return r, errors.Chain(err, "error unmarshaling json")
	}

	return r, nil
}

func (db DatabaseV2) DeleteRecord(id uuid.UUID) error {
	err := os.Remove(db.RecordPath(id))
	if err != nil {
		return errors.Chain(err, "error removing record file")
	}
	return nil
}
