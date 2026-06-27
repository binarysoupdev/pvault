package vault

import (
	"pvault/errors"
	"pvault/vault/index"
	"pvault/vault/record"
)

func (v Vault) updateIndex(r record.Record) error {
	existingName, ok := v.Index.FindName(r.ID)
	if ok && existingName != r.Name {
		delete(v.Index, existingName)
	}
	v.Index[r.Name] = r.ID

	err := index.Codec{}.Encode(v.IndexPath(), v.Index)
	if err != nil {
		return errors.Chain(err, "error saving index file")
	}

	return nil
}

func (v Vault) deleteIndex(name string) error {
	delete(v.Index, name)

	err := index.Codec{}.Encode(v.IndexPath(), v.Index)
	if err != nil {
		return errors.Chain(err, "error saving index file")
	}

	return nil
}
