package vault

import (
	"pvault/app/vault/database"

	"github.com/binarysoupdev/go-commando/errors"
)

func (v Vault) IndexPath() string {
	return v.DatabaseEncoder.IndexPath(v.Path)
}

func (v Vault) SaveIndex() error {
	err := database.SaveIndex(v.DatabaseEncoder, v.Path, v.Map)
	if err != nil {
		return errors.Chain(err, "error saving index map")
	}
	return nil
}

func (v *Vault) LoadIndex() error {
	var err error

	v.Map, err = database.LoadIndex(v.DatabaseEncoder, v.Path)
	if err != nil {
		return errors.Chain(err, "error loading index map")
	}
	return nil
}
