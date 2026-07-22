package vault

import (
	"path/filepath"
	"pvault/app/vault/database"
	"pvault/app/vault/index"
	"pvault/app/vault/meta"
)

const META_FILE = "VAULT"

type Vault struct {
	Meta     meta.Metadata
	Database database.Database
	Map      index.IndexMap
}

func (v Vault) GetVersion() int {
	return v.Meta.DatabaseVersion
}

func (v Vault) GetPath() string {
	return v.Meta.Path
}

func newMetadata(path string, dbVersion int) meta.Metadata {
	return meta.Metadata{
		Path:            filepath.Join(path, META_FILE),
		DatabaseVersion: dbVersion,
	}
}
