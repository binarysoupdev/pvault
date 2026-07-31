package vault

import (
	"pvault/app/vault/meta"

	"github.com/binarysoupdev/go-commando/errors"
)

func (v Vault) MetadataPath() string {
	return v.MetaEncoder.MetadataPath(v.Path)
}

func (v Vault) SaveMetadata() error {
	err := meta.SaveMetadata(v.MetaEncoder, v.Path, v.Meta)
	if err != nil {
		return errors.Chain(err, "error saving metadata")
	}
	return nil
}

func (v *Vault) LoadMetadata() error {
	var err error
	v.Meta, err = meta.LoadMetadata(v.MetaEncoder, v.Path)
	if err != nil {
		return errors.Chain(err, "error loading metadata")
	}
	return nil
}
