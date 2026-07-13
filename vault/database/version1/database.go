package version1

import (
	"os"
	"pvault/vault/database/version2"
	"pvault/vault/index"
	v1 "pvault/vault/record/version1"

	"github.com/binarysoupdev/go-commando/errors"
)

const VERSION = 1

type Database struct {
	Path string
}

func NewDatabase(path string) Database {
	return Database{
		Path: path,
	}
}

func (Database) GetVersion() uint16 {
	return VERSION
}

func (db Database) Initialize(idx index.IndexMap) error {
	err := db.SaveIndex(idx)
	if err != nil {
		return errors.Chain(err, "error saving index file")
	}

	return nil
}

func (db Database) Upgrade(idx index.IndexMap) error {
	target := version2.NewDatabase(db.Path)

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

		file.Write(v1.MarshalFromLegacy(name, raw))
		_ = os.Remove(legacyFile)
	}

	_ = os.Remove(db.IndexPath())
	return nil
}
