package vault

import (
	"pvault/vault/data"
	"pvault/vault/index"
)

const CURRENT_VERSION uint16 = 2

type Vault struct {
	Path     string
	Index    index.IndexMap
	Database data.Database
}

func (v Vault) Version() uint16 {
	return v.Database.Version()
}
