package v2_test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"pvault/vault/record"
	v2 "pvault/vault/record/encoder/legacy/v2"
	record_v1 "pvault/vault/record/record/legacy/v1"
	record_v2 "pvault/vault/record/record/v2"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeRecordReturnsErrorWhenRecordVersionUnsupported(t *testing.T) {
	//-- arrange
	e := v2.Encoder{}
	buffer := &bytes.Buffer{}

	RECORD := record.Mock{
		Version: record_v2.VERSION + 1,
	}

	//-- act
	res := e.EncodeRecord(buffer, "", RECORD)

	//-- assert
	assert.ErrorContains(t, res, fmt.Sprintf("unsupported record version \"%d\"", RECORD.Version))
}

func TestDecodeRecordReturnsErrorWhenErrorReadingHeader(t *testing.T) {
	//-- arrange
	e := v2.Encoder{}
	buffer := &bytes.Buffer{}

	//-- act
	_, res := e.DecodeRecord(buffer, "")

	//-- assert
	assert.ErrorContains(t, res, "error reading header")
}

func TestDecodeRecordReturnsErrorWhenRecordVersionUnsupported(t *testing.T) {
	//-- arrange
	e := v2.Encoder{}
	buffer := &bytes.Buffer{}

	const VERSION = record_v2.VERSION + 1

	HEADER := make([]byte, 2)
	binary.BigEndian.PutUint16(HEADER, VERSION)
	buffer.Write(HEADER)

	//-- act
	_, res := e.DecodeRecord(buffer, "")

	//-- assert
	assert.ErrorContains(t, res, fmt.Sprintf("unsupported record version \"%d\"", VERSION))
}

func TestEncodeDecodeRecordReturnsRecordV1AndNoError(t *testing.T) {
	//-- arrange
	e := v2.Encoder{}
	buffer := &bytes.Buffer{}

	const PASSWORD = "Password123!"

	RECORD := record_v1.Record{
		Name:          "name",
		Password:      "password",
		Username:      "username",
		URL:           "url",
		RecoveryCodes: []string{"code1", "code2"},
	}
	require.NoError(t, e.EncodeRecord(buffer, PASSWORD, RECORD))

	//-- act
	res, err := e.DecodeRecord(buffer, PASSWORD)
	require.NoError(t, err)

	//-- assert
	require.NoError(t, err)
	assert.Equal(t, RECORD, res)
}

func TestEncodeDecodeRecordReturnsRecordV2AndNoError(t *testing.T) {
	//-- arrange
	e := v2.Encoder{}
	buffer := &bytes.Buffer{}

	const PASSWORD = "Password123!"

	RECORD := record_v2.Record{
		ID:       uuid.New(),
		Name:     "name",
		Password: "password",
		Username: "username",
		Other:    map[string]any{"foo": "bar"},
	}
	require.NoError(t, e.EncodeRecord(buffer, PASSWORD, RECORD))

	//-- act
	res, err := e.DecodeRecord(buffer, PASSWORD)
	require.NoError(t, err)

	//-- assert
	require.NoError(t, err)
	assert.Equal(t, RECORD, res)
}
