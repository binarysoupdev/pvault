package v1_test

import (
	"bytes"
	"errors"
	"fmt"
	"pvault/app/vault/record"
	v1 "pvault/app/vault/record/encoder/v1"
	record_v1 "pvault/app/vault/record/record/v1"
	"pvault/util"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeRecordReturnsErrorWhenRecordVersionUnsupported(t *testing.T) {
	//-- arrange
	e := v1.Encoder{}
	buffer := &bytes.Buffer{}

	RECORD := record.Mock{
		Version: record_v1.VERSION + 1,
	}

	//-- act
	res := e.EncodeRecord(buffer, "", RECORD)

	//-- assert
	require.ErrorContains(t, res, fmt.Sprintf("unsupported record version \"%d\"", RECORD.Version))
}

func TestEncodeRecordReturnsErrorWhenErrorWritingData(t *testing.T) {
	//-- arrange
	e := v1.Encoder{}
	mock := &util.MockWriter{
		WriteErrors: []error{errors.New("")},
	}

	RECORD := record_v1.Record{}

	//-- act
	res := e.EncodeRecord(mock, "", RECORD)

	//-- assert
	require.ErrorContains(t, res, "error encoding record v1")
}

func TestDecodeRecordReturnsErrorWhenErrorDecodingHash(t *testing.T) {
	//-- arrange
	e := v1.Encoder{}
	buffer := &bytes.Buffer{}

	//-- act
	_, res := e.DecodeRecord(buffer, "")

	//-- assert
	require.ErrorContains(t, res, "error decoding record v1")
	require.ErrorContains(t, res, "error decoding hash prefix")
}

func TestDecodeRecordReturnsErrorWhenErrorDecryptingRecord(t *testing.T) {
	//-- arrange
	e := v1.Encoder{}
	buffer := &bytes.Buffer{}

	HASH := make([]byte, v1.HASH_SIZE)
	buffer.Write(HASH)

	//-- act
	_, res := e.DecodeRecord(buffer, "")

	//-- assert
	require.ErrorContains(t, res, "error decrypting record v1")
}

func TestEncodeDecodeRecordReturnsRecordAndNoError(t *testing.T) {
	//-- arrange
	e := v1.Encoder{}
	buffer := &bytes.Buffer{}

	const PASSWORD = "Password123!"
	RECORD := record_v1.Record{
		Password:      "password",
		Username:      "username",
		URL:           "url",
		RecoveryCodes: []string{"code1", "code2"},
	}

	err := e.EncodeRecord(buffer, PASSWORD, RECORD)
	require.NoError(t, err)

	//-- act
	res, err := e.DecodeRecord(buffer, PASSWORD)

	//-- assert
	require.NoError(t, err)
	assert.Equal(t, RECORD, res)
}
