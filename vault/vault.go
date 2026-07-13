package vault

import (
	"pvault/vault/database"
	"pvault/vault/index"
)

const CURRENT_VERSION uint16 = 2

type Vault struct {
	Path     string
	Index    index.IndexMap
	Database database.Database
}

func (v Vault) Version() uint16 {
	return v.Database.GetVersion()
}
