package data

import (
	"encoding/binary"
	"encoding/json"
	"os"
	"path/filepath"
	"pvault/errors"
	"pvault/vault/record"

	"github.com/binarysoupdev/cryptool/crypt"
	"github.com/google/uuid"
)

const CURRENT_RECORD_VERSION uint16 = 2

func (db DatabaseV2) RecordPath(id uuid.UUID) string {
	return filepath.Join(filepath.Dir(db.Path), id.String())
}

func (db DatabaseV2) SaveRecord(r record.Record, password string) error {
	c, salt := crypt.NewFromPassword(password)

	plaintext, err := json.Marshal(r)
	if err != nil {
		return errors.Chain(err, "error marshaling json")
	}

	ciphertext := c.Encrypt(plaintext)

	file, err := os.Create(db.RecordPath(r.ID))
	if err != nil {
		return errors.Chain(err, "error creating record file")
	}
	defer file.Close()

	version := make([]byte, 2)
	binary.BigEndian.PutUint16(version, CURRENT_RECORD_VERSION)

	file.Write(version)
	file.Write(salt)
	file.Write(ciphertext)

	return nil
}

func (db DatabaseV2) LoadRecord(id uuid.UUID, password string) (record.Record, error) {
	raw, err := os.ReadFile(db.RecordPath(id))
	if err != nil {
		return record.Record{}, errors.Chain(err, "error reading record file")
	}

	version := binary.BigEndian.Uint16(raw)
	raw = raw[2:]

	c := crypt.LoadFromPassword(password, raw[:crypt.SALT_SIZE])

	plaintext, err := c.Decrypt(raw[crypt.SALT_SIZE:])
	if err != nil {
		return record.Record{}, errors.Chain(err, "error decrypting ciphertext")
	}

	switch version {
	case 2:
		return db.parseRecordV2(plaintext)
	default:
		return record.Record{}, errors.Format("unsupported record version \"%d\"", version)
	}
}

func (DatabaseV2) parseRecordV2(raw []byte) (record.Record, error) {
	var r record.Record

	err := json.Unmarshal(raw, &r)
	if err != nil {
		return record.Record{}, errors.Chain(err, "error unmarshaling json")
	}

	return r, nil
}

func (db DatabaseV2) DeleteRecord(id uuid.UUID) error {
	err := os.Remove(db.RecordPath(id))
	if err != nil {
		return errors.Chain(err, "error removing record file")
	}
	return nil
}
