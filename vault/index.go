package vault

import (
	"pvault/chain"
	"pvault/vault/record"
)

func (v Vault) updateIndex(r record.Record) error {
	existingName, ok := v.Index.FindName(r.ID)
	if ok && existingName != r.Name {
		delete(v.Index, existingName)
	}
	v.Index[r.Name] = r.ID

	err := v.Index.Save(v.IndexPath())
	if err != nil {
		return chain.Error(err, "error saving index file")
	}

	return nil
}

func (v Vault) deleteIndex(name string) error {
	delete(v.Index, name)

	err := v.Index.Save(v.IndexPath())
	if err != nil {
		return chain.Error(err, "error saving index file")
	}

	return nil
}
