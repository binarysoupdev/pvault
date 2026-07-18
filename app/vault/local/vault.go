package local

import (
	"os"
	"pvault/app/vault/data"
	"pvault/app/vault/index"
	v2 "pvault/app/vault/index/version2"

	"github.com/binarysoupdev/go-commando/errors"
)

type Vault struct {
	Path  string
	Index index.Index
	Map   data.NameMap
}

func NewVault(path string) Vault {
	return Vault{
		Path: path,
	}
}

func (v Vault) GetVersion() int {
	return v.Index.GetVersion()
}

func (v Vault) IsOutOfDate() bool {
	return v.GetVersion() < v2.VERSION
}

func InitializeNewVault(path string) (Vault, error) {
	_, err := os.Stat(path)
	if err == nil || !os.IsNotExist(err) {
		return Vault{}, errors.New("vault path already exists")
	}

	err = os.MkdirAll(path, 0755)
	if err != nil {
		return Vault{}, errors.Chain(err, "error creating vault directory")
	}

	v := Vault{
		Path:  path,
		Index: v2.NewIndex(path),
		Map:   data.NameMap{},
	}

	err = v.Index.SaveMap(v.Map)
	if err != nil {
		return Vault{}, errors.Chain(err, "error saving initial index")
	}

	return v, nil
}
