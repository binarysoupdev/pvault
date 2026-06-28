package vault

import (
	"encoding/json"
	"os"
	"pvault/errors"
	"pvault/vault/record"

	"github.com/binarysoupdev/cryptool/crypt"
	"github.com/google/uuid"
)

func (v Vault) SaveRecord(r record.Record, password string) error {
	err := v.ValidateRecord(r)
	if err != nil {
		return errors.Chain(err, "error validating record")
	}

	err = v.saveEncryptedRecord(r, password)
	if err != nil {
		return err
	}

	err = v.updateIndex(r)
	if err != nil {
		return err
	}

	return nil
}

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

func (v Vault) updateIndex(r record.Record) error {
	existingName, ok := v.Index.FindName(r.ID)
	if ok && existingName != r.Name {
		delete(v.Index, existingName)
	}
	v.Index[r.Name] = r.ID

	err := v.Database.SaveIndex(v.Index)
	if err != nil {
		return errors.Chain(err, "error saving index file")
	}

	return nil
}

func (v Vault) LoadRecord(name string, password string) (record.Record, error) {
	id, ok := v.Index[name]
	if !ok {
		return record.Record{}, errors.Format("name \"%s\" not found", name)
	}

	r, err := v.loadEncryptedRecord(id, password)
	if err != nil {
		return record.Record{}, err
	}

	return r, nil
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

func (v Vault) DeleteRecord(name string) (uuid.UUID, error) {
	id, ok := v.Index[name]
	if !ok {
		return uuid.Nil, errors.Format("name \"%s\" not found", name)
	}

	err := os.Remove(v.RecordPath(id))
	if err != nil {
		return id, errors.Chain(err, "error deleting record file")
	}

	err = v.deleteIndex(name)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (v Vault) deleteIndex(name string) error {
	delete(v.Index, name)

	err := v.Database.SaveIndex(v.Index)
	if err != nil {
		return errors.Chain(err, "error saving index file")
	}

	return nil
}
