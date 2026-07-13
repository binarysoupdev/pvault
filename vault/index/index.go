package index

import "github.com/google/uuid"

type Index interface {
	GetVersion()
	RecordPath(id uuid.UUID) string

	SaveIndex(m IndexMap) error
	LoadIndex() (IndexMap, error)
	Upgrade(m IndexMap) error
}
