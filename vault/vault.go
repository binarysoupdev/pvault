package vault

import (
	"path/filepath"
	"pvault/cfg"
	"pvault/data"

	"github.com/google/uuid"
)

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

	return data.SaveJSON(r, filepath.Join(cfg.Global.VaultPath, r.ID.String()+".json"))
}

func (Vault) LoadRecord(id uuid.UUID) (Record, error) {
	return data.LoadJSON[Record](filepath.Join(cfg.Global.VaultPath, id.String()+".json"))
}
