package vault

import (
	"os"
	"pvault/errors"
	"pvault/vault/data/version2"
	"pvault/vault/index"
)

func InitializeNew(path string) (Vault, error) {
	_, err := os.Stat(path)
	if err == nil || !os.IsNotExist(err) {
		return Vault{}, errors.New("vault path already exists")
	}

	err = os.MkdirAll(path, 0755)
	if err != nil {
		return Vault{}, errors.Chain(err, "error creating vault directory")
	}

	v := Vault{
		Path:     path,
		Index:    index.IndexMap{},
		Database: version2.New(path),
	}

	err = v.Database.Initialize(v.Index)
	if err != nil {
		return Vault{}, errors.Chain(err, "error initializing vault database")
	}

	return v, nil
}

func (v *Vault) ReloadIndex() error {
	var err error

	v.Index, err = v.Database.LoadIndex()
	if err != nil {
		return errors.Chain(err, "error loading index from database")
	}

	return nil
}
