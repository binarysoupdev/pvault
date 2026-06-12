package vault

import (
	"path/filepath"
	"pvault/config"
	"pvault/data"

	"github.com/google/uuid"
)

func (Vault) SaveRecord(r Record) error {
	//TODO: check record name is unique

	return data.SaveJSON(r, filepath.Join(config.Global.VaultPath, r.ID.String()+".json"))
}

func (Vault) LoadRecord(id uuid.UUID) (Record, error) {
	return data.LoadJSON[Record](filepath.Join(config.Global.VaultPath, id.String()+".json"))
}
