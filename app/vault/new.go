package vault

import (
	"os"
	v3 "pvault/app/vault/database/database/v3"
	"pvault/app/vault/index"
	"pvault/app/vault/meta"

	"github.com/binarysoupdev/go-commando/errors"
)

func CreateNew(path, name string) (Vault, error) {
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
		Meta:     meta.New(v3.VERSION, name),
		Database: v3.Database{},
		Map:      index.IndexMap{},
	}

	err = v.SaveMetadata()
	if err != nil {
		return Vault{}, err
	}

	err = v.SaveIndex()
	if err != nil {
		return Vault{}, err
	}

	return v, nil
}
