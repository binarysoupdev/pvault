package index

import (
	"pvault/app/vault/data"
	v2 "pvault/app/vault/index/version2"

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
