package version1

import (
	"os"
	"pvault/errors"
	"pvault/vault/data"
	"pvault/vault/index"
)

func (db Database) Upgrade(idx index.IndexMap, target data.Database) error {
	for name, id := range idx {
		legacyFile := db.RecordPath(id)

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
