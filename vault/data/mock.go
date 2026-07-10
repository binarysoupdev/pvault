package data

import (
	"pvault/vault/index"
	"pvault/vault/record"

	"github.com/google/uuid"
)

type DatabaseMock struct {
	Version uint16
	Index   index.IndexMap
	Record  record.Record

	InitializeError   error
	UpgradeError      error
	SaveIndexError    error
	LoadIndexError    error
	SaveRecordError   error
	LoadRecordError   error
	DeleteRecordError error
}

func (db DatabaseMock) GetVersion() uint16 {
	return db.Version
}

func (db DatabaseMock) Initialize(idx index.IndexMap) error {
	return db.InitializeError
}

func (db *DatabaseMock) Upgrade(idx index.IndexMap, target Database) error {
	return db.UpgradeError
}

func (DatabaseMock) IndexPath() string {
	return ""
}

func (db *DatabaseMock) SaveIndex(idx index.IndexMap) error {
	db.Index = idx
	return db.SaveIndexError
}

func (db DatabaseMock) LoadIndex() (index.IndexMap, error) {
	return db.Index, db.LoadIndexError
}

func (DatabaseMock) RecordPath(id uuid.UUID) string {
	return ""
}

func (db *DatabaseMock) SaveRecord(r record.Record, password string) error {
	db.Record = r
	return db.SaveRecordError
}

func (db DatabaseMock) LoadRecord(id uuid.UUID, password string) (record.Record, error) {
	return db.Record, db.LoadRecordError
}

func (db *DatabaseMock) DeleteRecord(id uuid.UUID) error {
	db.Record = record.Record{}
	return db.DeleteRecordError
}
