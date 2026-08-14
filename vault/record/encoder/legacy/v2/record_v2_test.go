package v2_test

import (
	"bytes"

	v2 "pvault/vault/record/encoder/legacy/v2"
	record_v2 "pvault/vault/record/record/v2"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeV2ReturnsErrorWhenErrorDecryptingRecord(t *testing.T) {
	//-- arrange
	e := v2.Encoder{}
	buffer := &bytes.Buffer{}

	LENGTH := make([]byte, 2)
	buffer.Write(LENGTH)

	//-- act
	_, res := e.DecodeV2(buffer, "")

	//-- assert
	assert.ErrorContains(t, res, "error decrypting record v2")
}

func TestEncodeDecodeV2ReturnsRecordAndNoError(t *testing.T) {
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
	require.NoError(t, e.EncodeV2(buffer, PASSWORD, RECORD))

	//-- act
	res, err := e.DecodeV2(buffer, PASSWORD)
	require.NoError(t, err)

	//-- assert
	require.NoError(t, err)
	assert.Equal(t, RECORD, res)
}
