package vault

import (
	"pvault/chain"
	"pvault/data"

	"github.com/google/uuid"
)

func (v Vault) SaveRecord(r Record) error {
	//TODO: check record name is unique

	err := data.SaveJSON(r, v.RecordPath(r.ID))
	if err != nil {
		return chain.Error(err, "error saving record file")
	}

	v.Index[r.Name] = r.ID

	err = v.Index.Save(v.IndexPath())
	if err != nil {
		return chain.Error(err, "error saving index file")
	}

	return nil
}

func (v Vault) LoadRecord(id uuid.UUID) (Record, error) {
	r, err := data.LoadJSON[Record](v.RecordPath(id))
	if err != nil {
		return Record{}, chain.Error(err, "error loading record file")
	}

	return r, nil
}
