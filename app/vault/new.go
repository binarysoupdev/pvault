package vault

import (
	"os"
	db_v3 "pvault/app/vault/database/encoder/v3"
	"pvault/app/vault/index"
	"pvault/app/vault/meta"
	meta_v1 "pvault/app/vault/meta/encoder/v1"

	"github.com/binarysoupdev/go-commando/errors"
)

func New(path, nickname string) Vault {
	return Vault{
		Path:            path,
		MetaEncoder:     meta_v1.Encoder{},
		DatabaseEncoder: db_v3.Encoder{},
		Meta:            meta.New(db_v3.VERSION, nickname),
		Map:             index.IndexMap{},
	}
}

func InitializeNew(path, name string) (Vault, error) {
	_, err := os.Stat(path)
	if err == nil || !os.IsNotExist(err) {
		return Vault{}, errors.New("vault path already exists")
	}

	err = os.MkdirAll(path, 0755)
	if err != nil {
		return Vault{}, errors.Chain(err, "error creating vault directory")
	}

	v := New(path, name)

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
