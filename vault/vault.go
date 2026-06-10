package vault

import (
	"path/filepath"
	"pvault/config"
	"pvault/data"

	"github.com/google/uuid"
)

type Vault struct{}

func (v Vault) NewRecord(name string) (Record, error) {
	r := NewRecord(name)
	return r, v.SaveRecord(r)
}

func (Vault) SaveRecord(r Record) error {
	//TODO: check record name is unique

	return data.SaveJSON(r, filepath.Join(config.Global.VaultPath, r.ID.String()+".json"))
}

func (Vault) LoadRecord(id uuid.UUID) (Record, error) {
	return data.LoadJSON[Record](filepath.Join(config.Global.VaultPath, id.String()+".json"))
}
