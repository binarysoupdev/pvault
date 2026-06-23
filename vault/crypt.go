package vault

import (
	"pvault/chain"
	"pvault/data"
	"pvault/vault/record"

	"github.com/google/uuid"
)

func (v Vault) saveEncryptedRecord(r record.Record) error {
	err := data.SaveJSON(r, v.RecordPath(r.ID))
	if err != nil {
		return chain.Error(err, "error saving record file")
	}

	return nil
}

func (v Vault) loadEncryptedRecord(id uuid.UUID) (record.Record, error) {
	r, err := data.LoadJSON[record.Record](v.RecordPath(id))
	if err != nil {
		return record.Record{}, chain.Error(err, "error loading record file")
	}

	return r, nil
}
