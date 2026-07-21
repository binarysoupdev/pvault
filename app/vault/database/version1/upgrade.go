package v1

import (
	"encoding/binary"
	"os"
	v3 "pvault/app/vault/database/version3"
	"pvault/app/vault/index"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
)

func (db Database) Upgrade(idx index.IndexMap) (v3.Database, error) {
	target := v3.NewDatabase(db.Path)

	err := db.upgradeIndex(target, idx)
	if err != nil {
		return target, err
	}
	var errs errors.Errors

	for name, id := range idx {
		err := db.upgradeRecord(target, id, name)
		if err != nil {
			errs.Add(err)
			continue
		}
	}

	return target, errs.Collapse("\n")
}

func (db Database) upgradeIndex(target v3.Database, idx index.IndexMap) error {
	file, err := os.Create(target.IndexPath())
	if err != nil {
		return errors.Chain(err, "error creating new index file")
	}
	defer file.Close()

	version := make([]byte, 2)
	binary.BigEndian.PutUint16(version, uint16(target.GetVersion()))
	file.Write(version)

	err = target.EncodeIndex(file, idx)
	if err != nil {
		return errors.Chain(err, "error encoding index")
	}

	err = os.Remove(db.IndexPath())
	if err != nil {
		return errors.Chain(err, "error removing old index file")
	}

	return nil
}

func (db Database) upgradeRecord(target v3.Database, id uuid.UUID, name string) error {
	old, err := os.Open(db.RecordPath(id))
	if err != nil {
		return errors.Chain(err, "error opening old record file")
	}
	defer old.Close()

	data, err := db.DecodeRawV1(old)
	if err != nil {
		return errors.Chain(err, "error decoding old record")
	}

	new, err := os.Create(target.RecordPath(id))
	if err != nil {
		return errors.Chain(err, "error creating new record file")
	}
	defer new.Close()

	err = target.EncodeRawV1(new, data, id, name)
	if err != nil {
		return errors.Chain(err, "error encoding new record")
	}

	err = os.Remove(db.RecordPath(id))
	if err != nil {
		return errors.Chain(err, "error removing old record")
	}

	return nil
}
