package database

import (
	"io"
	"path/filepath"
	v3 "pvault/app/vault/database/encoder/v3"
	"pvault/app/vault/index"
	"pvault/app/vault/record"

	"github.com/google/uuid"
)

type EncoderMock struct {
	Version int

	Index            index.IndexMap
	EncodeIndexError error
	DecodeIndexError error

	Record            record.Record
	Password          string
	EncodeRecordError error
	DecodeRecordError error

	UpgradedEncoder v3.Encoder
	UpgradeError    error
}

func (m EncoderMock) GetVersion() int {
	return m.Version
}

func (EncoderMock) IndexPath(path string) string {
	return filepath.Join(path, "index.mock")
}

func (m *EncoderMock) EncodeIndex(_ io.Writer, idx index.IndexMap) error {
	m.Index = idx
	return m.EncodeIndexError
}

func (m EncoderMock) DecodeIndex(_ io.Reader) (index.IndexMap, error) {
	return m.Index, m.DecodeIndexError
}

func (EncoderMock) RecordPath(path string, id uuid.UUID) string {
	return filepath.Join(path, id.String()+".mock")
}

func (m *EncoderMock) EncodeRecord(_ io.Writer, password string, r record.Record) error {
	m.Password = password
	m.Record = r
	return m.EncodeRecordError
}

func (m EncoderMock) DecodeRecord(_ io.Reader, _ string) (record.Record, error) {
	return m.Record, m.DecodeRecordError
}

func (m EncoderMock) Upgrade(_ string) (v3.Encoder, error) {
	return m.UpgradedEncoder, m.UpgradeError
}
