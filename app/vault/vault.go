package vault

import (
	"pvault/app/vault/database"
	"pvault/app/vault/index"
	"pvault/app/vault/meta"
)

type Vault struct {
	Path            string
	MetaEncoder     meta.Encoder
	DatabaseEncoder database.Encoder

	Meta meta.Metadata
	Map  index.IndexMap
}

func (v Vault) GetVersion() int {
	return v.Meta.DatabaseVersion
}
