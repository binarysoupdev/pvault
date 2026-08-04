package v1

import (
	"os"
	"pvault/app/vault/database"
	v3 "pvault/app/vault/database/database/v3"
	"pvault/app/vault/index"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
)

func (db Database) Upgrade(path string) (v3.Database, error) {
	target := v3.Database{}

	idx, err := database.LoadIndex(db, path)
	if err != nil {
		return target, errors.Chain(err, "error loading old index file")
	}

	err = db.upgradeIndex(target, path, idx)
	if err != nil {
		return target, err
	}
	var errs errors.Errors

	for name, id := range idx {
		err := db.upgradeRecord(target, path, id, name)
		if err != nil {
			errs.Add(errors.Chain(err, "error backing record "+id.String()))
			continue
		}
	}

	return target, errs.Collapse("\n")
}

func (db Database) upgradeIndex(target v3.Database, path string, idx index.IndexMap) error {
	err := database.SaveIndex(target, path, idx)
	if err != nil {
		return errors.Chain(err, "error creating new index file")
	}

	err = os.Remove(db.IndexPath(path))
	if err != nil {
		return errors.Chain(err, "error removing old index file")
	}

	return nil
}

func (db Database) upgradeRecord(target v3.Database, path string, id uuid.UUID, name string) error {
	old, err := os.Open(db.RecordPath(path, id))
	if err != nil {
		return errors.Chain(err, "error opening old record file")
	}
	defer old.Close()

	data, err := db.DecodeRawV1(old)
	if err != nil {
		return errors.Chain(err, "error decoding old record")
	}

	new, err := os.Create(target.RecordPath(path, id))
	if err != nil {
		return errors.Chain(err, "error creating new record file")
	}
	defer new.Close()

	err = target.EncodeRawV1(new, data, id, name)
	if err != nil {
		return errors.Chain(err, "error encoding new record")
	}

	err = os.Remove(db.RecordPath(path, id))
	if err != nil {
		return errors.Chain(err, "error removing old record")
	}

	return nil
}
