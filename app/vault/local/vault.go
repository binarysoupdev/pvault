package local

import (
	"pvault/app/vault/data"
	"pvault/app/vault/index"
	v2 "pvault/app/vault/index/version2"
)

type Vault struct {
	Path  string
	Index index.Index
	Map   data.NameMap
}

func (v Vault) GetVersion() int {
	return v.Index.GetVersion()
}

func (v Vault) IsOutOfDate() bool {
	return v.GetVersion() < v2.VERSION
}
