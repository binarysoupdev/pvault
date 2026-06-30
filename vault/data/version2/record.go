package version2

import (
	"encoding/binary"
	"encoding/json"
	"os"
	"path/filepath"
	"pvault/errors"
	"pvault/vault/data"
	"pvault/vault/record"
	"pvault/vault/record/legacy"

	"github.com/binarysoupdev/cryptool/crypt"
	"github.com/google/uuid"
)

const RECORD_VERSION uint16 = 2

func (db Database) RecordPath(id uuid.UUID) string {
	return filepath.Join(db.Path, id.String())
}

func (db Database) SaveRecord(r record.Record, password string) error {
	return data.SaveEncryptedRecord(db.RecordPath(r.ID), password, RECORD_VERSION, nil, r)
}

func (db Database) LoadRecord(id uuid.UUID, password string) (record.Record, error) {
	raw, err := os.ReadFile(db.RecordPath(id))
	if err != nil {
		return record.Record{}, errors.Chain(err, "error reading record file")
	}

	version := binary.BigEndian.Uint16(raw)
	raw = raw[2:]

	switch version {
	case 1:
		return db.parseRecordV1(password, raw, id)
	case 2:
		return db.parseRecordV2(password, raw)
	default:
		return record.Record{}, errors.Format("unsupported record version \"%d\"", version)
	}
}

func (Database) parseRecordV1(password string, raw []byte, id uuid.UUID) (record.Record, error) {
	length := binary.BigEndian.Uint16(raw)
	raw = raw[2:]

	name := string(raw[:length])
	raw = raw[length:]

	r, err := decryptJSON[legacy.RecordV1](password, raw)
	if err != nil {
		return record.Record{}, err
	}

	return r.Upgrade(id, name), nil
}

func (Database) parseRecordV2(password string, raw []byte) (record.Record, error) {
	return decryptJSON[record.Record](password, raw)
}

func decryptJSON[T any](password string, ciphertext []byte) (T, error) {
	var obj T
	c := crypt.LoadFromPassword(password, ciphertext[:crypt.SALT_SIZE])

	plaintext, err := c.Decrypt(ciphertext[crypt.SALT_SIZE:])
	if err != nil {
		return obj, errors.Chain(err, "error decrypting ciphertext")
	}

	err = json.Unmarshal(plaintext, &obj)
	if err != nil {
		return obj, errors.Chain(err, "error unmarshaling json")
	}

	return obj, nil
}

func (db Database) DeleteRecord(id uuid.UUID) error {
	err := os.Remove(db.RecordPath(id))
	if err != nil {
		return errors.Chain(err, "error removing record file")
	}
	return nil
}
