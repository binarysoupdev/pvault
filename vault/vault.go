package vault

import (
	"pvault/vault/database"
	"pvault/vault/index"
)

const CURRENT_VERSION int = 2

type Vault struct {
	Path     string
	Index    index.IndexMap
	Database database.Database
}

func (v Vault) Version() int {
	return v.Database.GetVersion()
}
