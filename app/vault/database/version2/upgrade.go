package v2

import (
	"encoding/binary"
	"io"
	"os"
	v3 "pvault/app/vault/database/version3"
	"pvault/app/vault/index"
	record_v1 "pvault/app/vault/record/version1"
	record_v2 "pvault/app/vault/record/version2"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
)

func (db Database) Upgrade(idx index.IndexMap) (v3.Database, error) {
	target := v3.NewDatabase(db.Path)

	err := db.upgradeIndex(target)
	if err != nil {
		return target, err
	}
	var errs errors.Errors

	for _, id := range idx {
		err := db.upgradeRecord(target, id)
		if err != nil {
			errs.Add(err)
			continue
		}
	}

	return target, errs.Collapse("\n")
}

func (db Database) upgradeIndex(target v3.Database) error {
	err := os.Rename(db.IndexPath(), target.IndexPath())
	if err != nil {
		return errors.Chain(err, "error renaming index file")
	}

	file, err := os.OpenFile(target.IndexPath(), os.O_WRONLY, 0)
	if err != nil {
		return errors.Chain(err, "error opening renamed index file")
	}
	defer file.Close()

	version := make([]byte, 2)
	binary.BigEndian.PutUint16(version, uint16(target.GetVersion()))
	file.Write(version)

	return nil
}

func (db Database) upgradeRecord(target v3.Database, id uuid.UUID) error {
	old, err := os.Open(db.RecordPath(id))
	if err != nil {
		return errors.Chain(err, "error opening old record file")
	}
	defer old.Close()

	header := make([]byte, 2)
	old.Read(header)

	version := binary.BigEndian.Uint16(header)

	switch version {
	case record_v1.VERSION:
		return db.upgradeRecordV1(target, old, id)
	case record_v2.VERSION:
		return nil
	default:
		return errors.Format("unsupported record version \"%d\"", version)
	}
}

func (db Database) upgradeRecordV1(target v3.Database, r io.Reader, id uuid.UUID) error {
	raw, err := db.DecodeRawV1(r)
	if err != nil {
		return errors.Chain(err, "error decoding old record")
	}

	new, err := os.Create(target.RecordPath(id))
	if err != nil {
		return errors.Chain(err, "error creating new record file")
	}
	defer new.Close()

	err = target.EncodeRawV1(new, raw.Data, id, raw.Name)
	if err != nil {
		return errors.Chain(err, "error encoding new record")
	}

	err = os.Remove(db.RecordPath(id))
	if err != nil {
		return errors.Chain(err, "error removing old record")
	}

	return nil
}
