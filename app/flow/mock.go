package flow

import (
	"pvault/app/vault/record"

	"github.com/google/uuid"
)

type Mock struct {
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

func (m Mock) GetVersion() int {
	return m.Version
}

func (m *Mock) SearchNames(term string) []string {
	m.SearchTermParam = term
	return m.SearchResults
}

func (m *Mock) ValidateRecord(r record.Record) error {
	m.Record = r
	return m.ValidateRecordError
}

func (m *Mock) SaveRecord(r record.Record, password string) error {
	m.Record = r
	m.PasswordParam = password
	return m.SaveRecordError
}

func (m *Mock) LoadRecord(name string, password string) (record.Record, error) {
	m.NameParam = name
	m.PasswordParam = password
	return m.Record, m.LoadRecordError
}

func (m *Mock) DeleteRecord(name string) (uuid.UUID, error) {
	m.NameParam = name
	return m.Record.GetID(), m.DeleteRecordError
}
