package version1

import (
	"os"
	"pvault/errors"
	"pvault/vault/data"
	"pvault/vault/index"
	"pvault/vault/record/legacy"

	"github.com/google/uuid"
)

const LEGACY_HASH_SIZE = 60

func (db Database) Upgrade(idx index.IndexMap, target data.Database) error {
	for name, id := range idx {
		legacyFile := db.LegacyRecordPath(id)

		raw, err := os.ReadFile(legacyFile)
		if err != nil {
			continue
		}

		file, err := os.Create(target.RecordPath(id))
		if err != nil {
			return errors.Chain(err, "error creating converted record file")
		}
		defer file.Close()

		file.Write(db.buildRecordV1Header(name))
		file.Write(raw[LEGACY_HASH_SIZE:])

		_ = os.Remove(legacyFile)
	}

	_ = os.Remove(db.IndexPath())
	return nil
}

func (db Database) LegacyRecordPath(id uuid.UUID) string {
	return db.RecordPath(id) + ".crypt"
}

func (db Database) SaveLegacyRecord(id uuid.UUID, password string, r legacy.RecordV1) error {
	hash := make([]byte, LEGACY_HASH_SIZE)
	return data.SaveEncryptedRecord(db.LegacyRecordPath(id), password, hash, r)
}
