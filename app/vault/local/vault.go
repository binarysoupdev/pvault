package local

import (
	"pvault/app/vault/data"
	"pvault/app/vault/index"
	index_v2 "pvault/app/vault/index/version2"
)

const CURRENT_VERSION = index_v2.VERSION

type Vault struct {
	Path  string
	Index index.Index
	Map   data.NameMap
}

func (v Vault) GetVersion() int {
	return v.Index.GetVersion()
}
