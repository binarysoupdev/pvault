package vault

import (
	"path/filepath"

	"github.com/google/uuid"
)

const VAULT_DIR = "tmp/vault"

type Vault struct{}

func (v Vault) NewRecord(name string) (Record, error) {
	r := Record{
		ID:       uuid.New(),
		Name:     name,
		Username: "",
		Password: "",
		Other:    map[string]string{},
	}

	return r, v.SaveRecord(r)
}

func (Vault) SaveRecord(r Record) error {
	//TODO: check record name is unique

	return SaveRecordJSON(r, filepath.Join(VAULT_DIR, r.ID.String()+".json"))
}

func (Vault) LoadRecord(id uuid.UUID) (Record, error) {
	return LoadRecordJSON(filepath.Join(VAULT_DIR, id.String()+".json"))
}
