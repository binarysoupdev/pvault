package remote

import (
	"errors"
	"pvault/vault/record"

	"github.com/google/uuid"
)

func (v Vault) ValidateRecord(r record.Record) error {
	return errors.New("not implemented")
}

func (v Vault) SaveRecord(r record.Record, password string) error {
	return errors.New("not implemented")
}

func (v Vault) LoadRecord(name string, password string) (record.Record, error) {
	return nil, errors.New("not implemented")
}

func (v Vault) DeleteRecord(name string) (uuid.UUID, error) {
	return uuid.Nil, errors.New("not implemented")
}
