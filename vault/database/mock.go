package database

import (
	"pvault/vault/index"
	v2 "pvault/vault/record/version/v2"

	"github.com/google/uuid"
)

type DatabaseMock struct {
	Version uint16
	Index   index.IndexMap
	Record  v2.Record

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

func (db *DatabaseMock) Upgrade(idx index.IndexMap) error {
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

func (db *DatabaseMock) SaveRecord(r v2.Record, password string) error {
	db.Record = r
	return db.SaveRecordError
}

func (db DatabaseMock) LoadRecord(id uuid.UUID, password string) (v2.Record, error) {
	return db.Record, db.LoadRecordError
}

func (db *DatabaseMock) DeleteRecord(id uuid.UUID) error {
	db.Record = v2.Record{}
	return db.DeleteRecordError
}
