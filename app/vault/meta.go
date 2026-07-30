package vault

import (
	"path/filepath"
	"pvault/app/vault/meta"
	meta_v1 "pvault/app/vault/meta/encoder/v1"

	"github.com/binarysoupdev/go-commando/errors"
)

const METADATA_FILE = "VAULT"

func (v Vault) MetadataPath() string {
	return filepath.Join(v.Path, METADATA_FILE)
}

func (v Vault) SaveMetadata() error {
	err := meta.SaveMetadata(meta_v1.Encoder{}, v.MetadataPath(), v.Meta)
	if err != nil {
		return errors.Chain(err, "error saving metadata")
	}
	return nil
}

func (v *Vault) LoadMetadata() error {
	var err error
	v.Meta, err = meta.LoadMetadata(meta_v1.Encoder{}, v.MetadataPath())
	if err != nil {
		return errors.Chain(err, "error loading metadata")
	}
	return nil
}
