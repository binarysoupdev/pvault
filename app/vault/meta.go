package vault

import (
	"path/filepath"
	"pvault/app/vault/meta"
)

const METADATA_FILE = "VAULT"

func createNewMetadata(dbVersion int) meta.Metadata {
	return meta.New(dbVersion)
}

func metadataPath(path string) string {
	return filepath.Join(path, METADATA_FILE)
}
