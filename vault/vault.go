package vault

import (
	"os"
	"pvault/errors"
	"pvault/vault/index"
)

const VERSION = index.VERSION

type Vault struct {
	Path  string
	Index index.IndexMap
}

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

	err = v.Index.Save(path)
	if err != nil {
		return v, errors.Chain(err, "error saving index file")
	}

	return v, nil
}

func Open(path string) (Vault, error) {
	v := Vault{
		Path: path,
	}

	var err error
	v.Index, err = index.LoadIndex(v.Path)
	if err != nil {
		return v, errors.Chain(err, "error loading index file")
	}

	return v, nil
}
