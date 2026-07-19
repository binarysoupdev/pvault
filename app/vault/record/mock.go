package record

import (
	"os"
	v2 "pvault/app/vault/record/version2"

	"github.com/google/uuid"
)

type Mock struct {
	Version       int
	ID            uuid.UUID
	Name          string
	UpgradeRecord v2.Record

	PasswordParam string
	SaveFileError error
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

func (m *Mock) SaveFile(path string, password string) error {
	m.PasswordParam = password
	os.WriteFile(path, []byte{}, 0666)
	return m.SaveFileError
}

func (m Mock) Upgrade() v2.Record {
	return m.UpgradeRecord
}
