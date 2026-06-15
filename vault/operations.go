package vault

import (
	"fmt"
	"pvault/chain"
	"pvault/data"

	"github.com/google/uuid"
)

func (v Vault) SaveRecord(r Record) error {
	existingId, ok := v.Index[r.Name]
	if ok && existingId != r.ID {
		return fmt.Errorf("name \"%s\" already exists", r.Name)
	}

	err := data.SaveJSON(r, v.RecordPath(r.ID))
	if err != nil {
		return chain.Error(err, "error saving record file")
	}

	existingName, ok := v.Index.findName(r.ID)
	if ok && existingName != r.Name {
		delete(v.Index, existingName)
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
