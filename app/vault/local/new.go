package local

import (
	"os"
	"pvault/app/vault/data"
	"pvault/app/vault/index"
	v2 "pvault/app/vault/index/version2"

	"github.com/binarysoupdev/go-commando/errors"
)

func CreateNewVault(path string) (Vault, error) {
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

func OpenVault(path string) (Vault, error) {
	var err error
	v := Vault{
		Path: path,
	}

	v.Index, err = index.Load(v.Path)
	if err != nil {
		return Vault{}, errors.Chain(err, "error finding index")
	}

	v.Map, err = v.Index.LoadMap()
	if err != nil {
		return Vault{}, errors.Chain(err, "error loading index")
	}

	return v, nil
}
