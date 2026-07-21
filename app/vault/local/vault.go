package local

import (
	"pvault/app/vault/database"
	db_v3 "pvault/app/vault/database/version3"
	"pvault/app/vault/index"
)

const CURRENT_VERSION = db_v3.VERSION

type Vault struct {
	Database database.Database
	Map      index.IndexMap
}

func (v Vault) GetVersion() int {
	return v.Database.GetVersion()
}

func (v Vault) GetPath() string {
	return v.Database.GetPath()
}
