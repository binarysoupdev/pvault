package v1

import (
	"os"
	"pvault/vault/database"
	v3 "pvault/vault/database/encoder/v3"
	"pvault/vault/index"

	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/google/uuid"
)

func (db Encoder) Upgrade(path string) (v3.Encoder, error) {
	target := v3.Encoder{}

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

func (db Encoder) upgradeIndex(target v3.Encoder, path string, idx index.IndexMap) error {
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

func (db Encoder) upgradeRecord(target v3.Encoder, path string, id uuid.UUID, name string) error {
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
