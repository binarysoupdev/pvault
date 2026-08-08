package v2

import (
	"encoding/binary"
	"io"
	"os"
	"pvault/vault/database"
	v3 "pvault/vault/database/encoder/v3"
	record_v1 "pvault/vault/record/record/legacy/v1"
	record_v2 "pvault/vault/record/record/v2"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
)

func (db Encoder) Upgrade(path string) (v3.Encoder, error) {
	target := v3.Encoder{}

	idx, err := database.LoadIndex(db, path)
	if err != nil {
		return target, errors.Chain(err, "error loading old index file")
	}

	err = db.upgradeIndex(target, path)
	if err != nil {
		return target, err
	}
	var errs errors.Errors

	for _, id := range idx {
		err := db.upgradeRecord(target, path, id)
		if err != nil {
			errs.Add(errors.Chain(err, "error backing record "+id.String()))
			continue
		}
	}

	return target, errs.Collapse("\n")
}

func (db Encoder) upgradeIndex(target v3.Encoder, path string) error {
	err := os.Rename(db.IndexPath(path), target.IndexPath(path))
	if err != nil {
		return errors.Chain(err, "error renaming index file")
	}

	return nil
}

func (db Encoder) upgradeRecord(target v3.Encoder, path string, id uuid.UUID) error {
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

func (db Encoder) upgradeRecordV1(target v3.Encoder, r io.Reader, path string, id uuid.UUID) error {
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

	return nil
}
