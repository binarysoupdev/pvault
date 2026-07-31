package database

import (
	"io"
	"path/filepath"
	v3 "pvault/app/vault/database/database/v3"
	"pvault/app/vault/index"
	"pvault/app/vault/record"

	"github.com/google/uuid"
)

type Mock struct {
	Version int

	Index            index.IndexMap
	EncodeIndexError error
	DecodeIndexError error

	Record            record.Record
	Password          string
	EncodeRecordError error
	DecodeRecordError error

	UpgradedDatabase v3.Database
	UpgradeError     error
}

func (m Mock) GetVersion() int {
	return m.Version
}

func (Mock) IndexPath(path string) string {
	return filepath.Join(path, "index.mock")
}

func (m *Mock) EncodeIndex(_ io.Writer, idx index.IndexMap) error {
	m.Index = idx
	return m.EncodeIndexError
}

func (m Mock) DecodeIndex(_ io.Reader) (index.IndexMap, error) {
	return m.Index, m.DecodeIndexError
}

func (Mock) RecordPath(path string, id uuid.UUID) string {
	return filepath.Join(path, id.String()+".mock")
}

func (m *Mock) EncodeRecord(_ io.Writer, password string, r record.Record) error {
	m.Password = password
	m.Record = r
	return m.EncodeRecordError
}

func (m Mock) DecodeRecord(_ io.Reader, _ string) (record.Record, error) {
	return m.Record, m.DecodeRecordError
}

func (m Mock) Upgrade(_ string) (v3.Database, error) {
	return m.UpgradedDatabase, m.UpgradeError
}
