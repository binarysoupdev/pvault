package local

import (
	"pvault/app/vault/index"

	"github.com/binarysoupdev/go-commando/errors"
)

func (v *Vault) Load() error {
	var err error

	v.Index, err = index.Load(v.Path)
	if err != nil {
		return errors.Chain(err, "error finding index")
	}

	v.Map, err = v.Index.LoadMap()
	if err != nil {
		return errors.Chain(err, "error loading index")
	}

	return nil
}
