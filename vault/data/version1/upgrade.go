package version1

import (
	"os"
	"pvault/errors"
	"pvault/vault/data/version2"
	"pvault/vault/index"
)

func (db Database) UpgradeToVersion2(idx index.IndexMap) error {
	const LEGACY_HASH_SIZE = 60
	target := version2.NewDatabase(db.Path)

	for name, id := range idx {
		oldFile := db.RecordPath(id)

		raw, err := os.ReadFile(oldFile)
		if err != nil {
			continue
		}

		file, err := os.Create(target.RecordPath(id))
		if err != nil {
			return errors.Chain(err, "error creating converted record file")
		}
		defer file.Close()

		db.writeRecordMeta(file, name)
		file.Write(raw[LEGACY_HASH_SIZE:])

		_ = os.Remove(oldFile)
	}

	_ = os.Remove(db.IndexPath())
	return nil
}
