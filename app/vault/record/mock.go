package record

import (
	v2 "pvault/app/vault/record/record/v2"

	"github.com/google/uuid"
)

type Mock struct {
	Version       int
	ID            uuid.UUID
	Name          string
	UpgradeRecord v2.Record
}

func (m Mock) GetVersion() int {
	return m.Version
}

func (m Mock) GetID() uuid.UUID {
	return m.ID
}

func (m Mock) GetName() string {
	return m.Name
}

func (m Mock) Upgrade() v2.Record {
	return m.UpgradeRecord
}
