package version1

import (
	"encoding/binary"
	"os"
	"pvault/vault/data"
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

		ciphertext := raw[LEGACY_HASH_SIZE:]

		bytes := make([]byte, 2+len(name)+len(ciphertext))
		binary.BigEndian.PutUint16(bytes, uint16(len(name)))
		copy(bytes[2:], []byte(name))
		copy(bytes[2+len(name):], ciphertext)

		err = data.SaveVersionedRecord(target.RecordPath(id), RECORD_VERSION, bytes)
		if err != nil {
			return err
		}

		_ = os.Remove(oldFile)
	}

	_ = os.Remove(db.IndexPath())
	return nil
}
