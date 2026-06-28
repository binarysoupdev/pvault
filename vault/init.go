package vault

import (
	"os"
	"pvault/errors"
	"pvault/vault/data"
	"pvault/vault/index"
)

func InitializeNew(path string) (Vault, error) {
	v := Vault{
		Path:  path,
		Index: index.IndexMap{},
	}

	_, err := os.Stat(path)
	if err == nil || !os.IsNotExist(err) {
		return v, errors.New("vault path already exists")
	}

	err = os.MkdirAll(path, 0755)
	if err != nil {
		return v, errors.Chain(err, "error creating vault directory")
	}

	err = data.CurrentDatabase{}.SaveIndex(v.IndexPath(), v.Index)
	if err != nil {
		return v, errors.Chain(err, "error saving index file")
	}

	return v, nil
}
