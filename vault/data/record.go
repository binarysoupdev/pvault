package data

import (
	"encoding/binary"
	"encoding/json"
	"os"
	"path/filepath"
	"pvault/errors"
	"pvault/vault/record"
	"pvault/vault/record/legacy"

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

	switch version {
	case 1:
		return db.parseRecordV1(password, raw, id)
	case 2:
		return db.parseRecordV2(password, raw)
	default:
		return record.Record{}, errors.Format("unsupported record version \"%d\"", version)
	}
}

func (DatabaseV2) parseRecordV1(password string, raw []byte, id uuid.UUID) (record.Record, error) {
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

func (DatabaseV2) parseRecordV2(password string, raw []byte) (record.Record, error) {
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

func (db DatabaseV2) DeleteRecord(id uuid.UUID) error {
	err := os.Remove(db.RecordPath(id))
	if err != nil {
		return errors.Chain(err, "error removing record file")
	}
	return nil
}
