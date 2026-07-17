package vault

import (
	"pvault/app/vault/index"

	"github.com/binarysoupdev/go-commando/errors"
)

func Load(path string) (Vault, error) {
	v := Vault{
		Path: path,
	}
	var err error

	v.Index, err = index.Load(path)
	if err != nil {
		return Vault{}, errors.Chain(err, "error finding index")
	}

	v.Map, err = v.Index.LoadMap()
	if err != nil {
		return Vault{}, errors.Chain(err, "error loading index")
	}

	return v, nil
}
