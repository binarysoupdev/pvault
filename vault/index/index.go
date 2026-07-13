package index

import "github.com/google/uuid"

type Index interface {
	GetVersion()
	RecordPath(id uuid.UUID) string

	Save(m IndexMap) error
	Load() (IndexMap, error)
	Upgrade(m IndexMap) error
}
