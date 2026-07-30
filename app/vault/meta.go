package vault

import (
	"path/filepath"
	"pvault/app/vault/meta"
	meta_v1 "pvault/app/vault/meta/encoder/v1"

	"github.com/binarysoupdev/go-commando/errors"
)

const METADATA_FILE = "VAULT"

func metadataPath(path string) string {
	return filepath.Join(path, METADATA_FILE)
}

func saveMetadata(path string, m meta.Metadata) error {
	err := meta.SaveMetadata(meta_v1.Encoder{}, metadataPath(path), m)
	if err != nil {
		return errors.Chain(err, "error saving metadata")
	}
	return nil
}

func loadMetadata(path string) (meta.Metadata, error) {
	m, err := meta.LoadMetadata(meta_v1.Encoder{}, metadataPath(path))
	if err != nil {
		return meta.Metadata{}, errors.Chain(err, "error loading metadata")
	}
	return m, nil
}
