package index

import (
	"pvault/vault/data"
	v2 "pvault/vault/index/version2"

	"github.com/google/uuid"
)

type Index interface {
	GetVersion() int

	Filepath() string
	RecordPath(id uuid.UUID) string

	SaveMap(m data.NameMap) error
	LoadMap() (data.NameMap, error)

	Upgrade() (v2.Index, error)
}
