package version1

import (
	"path/filepath"
	"pvault/common"
	"pvault/vault/record"

	"github.com/google/uuid"
)

func (db Database) RecordPath(id uuid.UUID) string {
	return filepath.Join(db.Path, id.String()+".crypt")
}

func (db Database) SaveRecord(r record.Record, password string) error {
	return common.NotSupportedError{}
}

func (Database) LoadRecord(id uuid.UUID, password string) (record.Record, error) {
	return nil, common.NotSupportedError{}
}

func (Database) DeleteRecord(id uuid.UUID) error {
	return common.NotSupportedError{}
}
