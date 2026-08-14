package v2_test

import (
	"bytes"
	"encoding/binary"

	v2 "pvault/vault/record/encoder/legacy/v2"
	record_v1 "pvault/vault/record/record/legacy/v1"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeV1ReturnsErrorWhenErrorReadingLength(t *testing.T) {
	//-- arrange
	e := v2.Encoder{}
	buffer := &bytes.Buffer{}

	//-- act
	_, res := e.DecodeV1(buffer, "", uuid.Nil)

	//-- assert
	assert.ErrorContains(t, res, "error decoding record v1")
	assert.ErrorContains(t, res, "error reading length")
}

func TestDecodeV1ReturnsErrorWhenErrorReadingName(t *testing.T) {
	//-- arrange
	e := v2.Encoder{}
	buffer := &bytes.Buffer{}

	LENGTH := make([]byte, 2)
	binary.BigEndian.PutUint16(LENGTH, 1)
	buffer.Write(LENGTH)

	//-- act
	_, res := e.DecodeV1(buffer, "", uuid.Nil)

	//-- assert
	assert.ErrorContains(t, res, "error decoding record v1")
	assert.ErrorContains(t, res, "error reading name")
}

func TestDecodeV1ReturnsErrorWhenErrorDecryptingRecord(t *testing.T) {
	//-- arrange
	e := v2.Encoder{}
	buffer := &bytes.Buffer{}

	LENGTH := make([]byte, 2)
	buffer.Write(LENGTH)

	//-- act
	_, res := e.DecodeV1(buffer, "", uuid.Nil)

	//-- assert
	assert.ErrorContains(t, res, "error decrypting record v1")
}

func TestEncodeDecodeV1ReturnsRecordAndNoError(t *testing.T) {
	//-- arrange
	e := v2.Encoder{}
	buffer := &bytes.Buffer{}

	const PASSWORD = "Password123!"

	RECORD := record_v1.Record{
		ID:            uuid.New(),
		Name:          "name",
		Password:      "password",
		Username:      "username",
		URL:           "url",
		RecoveryCodes: []string{"code1", "code2"},
	}
	require.NoError(t, e.EncodeV1(buffer, PASSWORD, RECORD))

	//-- act
	res, err := e.DecodeV1(buffer, PASSWORD, RECORD.ID)
	require.NoError(t, err)

	//-- assert
	require.NoError(t, err)
	assert.Equal(t, RECORD, res)
}
