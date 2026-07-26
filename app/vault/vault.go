package vault

import (
	"pvault/app/vault/database"
	"pvault/app/vault/index"
	"pvault/app/vault/meta"
)

type Vault struct {
	Path     string
	Meta     meta.Metadata
	Database database.Database
	Map      index.IndexMap
}

func (v Vault) GetVersion() int {
	return v.Meta.DatabaseVersion
}
