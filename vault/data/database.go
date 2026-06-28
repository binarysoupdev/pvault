package data

import "pvault/vault/index"

type Database interface {
	Version() uint16
	SaveIndex(idx index.IndexMap) error
	LoadIndex() (index.IndexMap, error)
	Upgrade() error
}
