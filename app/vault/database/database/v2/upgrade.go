package v2

import (
	"encoding/binary"
	"io"
	"os"
	v3 "pvault/app/vault/database/database/v3"
	"pvault/app/vault/index"
	record_v1 "pvault/app/vault/record/record/v1"
	record_v2 "pvault/app/vault/record/record/v2"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
)

func (db Database) Upgrade(path string, idx index.IndexMap) (v3.Database, error) {
	target := v3.Database{}

	err := db.upgradeIndex(target, path)
	if err != nil {
		return target, err
	}
	var errs errors.Errors

	for _, id := range idx {
		err := db.upgradeRecord(target, path, id)
		if err != nil {
			errs.Add(err)
			continue
		}
	}

	return target, errs.Collapse("\n")
}

func (db Database) upgradeIndex(target v3.Database, path string) error {
	err := os.Rename(db.IndexPath(path), target.IndexPath(path))
	if err != nil {
		return errors.Chain(err, "error renaming index file")
	}

	file, err := os.OpenFile(target.IndexPath(path), os.O_WRONLY, 0)
	if err != nil {
		return errors.Chain(err, "error opening renamed index file")
	}
	defer file.Close()

	version := make([]byte, 2)
	binary.BigEndian.PutUint16(version, uint16(index.VERSION))
	file.Write(version)

	return nil
}

func (db Database) upgradeRecord(target v3.Database, path string, id uuid.UUID) error {
	old, err := os.Open(db.RecordPath(path, id))
	if err != nil {
		return errors.Chain(err, "error opening old record file")
	}
	defer old.Close()

	header := make([]byte, 2)
	old.Read(header)

	version := binary.BigEndian.Uint16(header)

	switch version {
	case record_v1.VERSION:
		return db.upgradeRecordV1(target, old, path, id)
	case record_v2.VERSION:
		return nil
	default:
		return errors.Format("unsupported record version \"%d\"", version)
	}
}

func (db Database) upgradeRecordV1(target v3.Database, r io.Reader, path string, id uuid.UUID) error {
	raw, err := db.DecodeRawV1(r)
	if err != nil {
		return errors.Chain(err, "error decoding old record")
	}

	new, err := os.Create(target.RecordPath(path, id))
	if err != nil {
		return errors.Chain(err, "error creating new record file")
	}
	defer new.Close()

	err = target.EncodeRawV1(new, raw.Data, id, raw.Name)
	if err != nil {
		return errors.Chain(err, "error encoding new record")
	}

	err = os.Remove(db.RecordPath(path, id))
	if err != nil {
		return errors.Chain(err, "error removing old record")
	}

	return nil
}
