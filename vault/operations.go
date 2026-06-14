package vault

import (
	"path/filepath"
	"pvault/chain"
	"pvault/data"

	"github.com/google/uuid"
)

func (v Vault) SaveRecord(r Record) error {
	//TODO: check record name is unique

	v.Index[r.Name] = r.ID

	err := v.Index.Save(v.indexPath())
	if err != nil {
		return chain.Error(err, "error saving index file")
	}

	err = data.SaveJSON(r, filepath.Join(v.Path, r.ID.String()+".json"))
	if err != nil {
		return chain.Error(err, "error saving record file")
	}

	return nil
}

func (v Vault) LoadRecord(id uuid.UUID) (Record, error) {
	return data.LoadJSON[Record](filepath.Join(v.Path, id.String()+".json"))
}
