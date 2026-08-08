package vault_flow

import (
	"pvault/app/vault/record"

	"github.com/google/uuid"
)

type VaultMock struct {
	Version int

	Record        record.Record
	SearchResults []string

	PasswordParam   string
	NameParam       string
	SearchTermParam string

	ValidateRecordError error
	SaveRecordError     error
	LoadRecordError     error
	DeleteRecordError   error
}

func (m VaultMock) GetVersion() int {
	return m.Version
}

func (m *VaultMock) SearchNames(term string) []string {
	m.SearchTermParam = term
	return m.SearchResults
}

func (m *VaultMock) ValidateRecord(r record.Record) error {
	m.Record = r
	return m.ValidateRecordError
}

func (m *VaultMock) SaveRecord(r record.Record, password string) error {
	m.Record = r
	m.PasswordParam = password
	return m.SaveRecordError
}

func (m *VaultMock) LoadRecord(name string, password string) (record.Record, error) {
	m.NameParam = name
	m.PasswordParam = password
	return m.Record, m.LoadRecordError
}

func (m *VaultMock) DeleteRecord(name string) (uuid.UUID, error) {
	m.NameParam = name
	return m.Record.GetID(), m.DeleteRecordError
}
