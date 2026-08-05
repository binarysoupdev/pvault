package v1_test

import (
	"bytes"
	"fmt"
	"pvault/app/vault/record"
	v1 "pvault/app/vault/record/encoder/legacy/v1"
	record_v1 "pvault/app/vault/record/record/legacy/v1"
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
	assert.ErrorContains(t, res, fmt.Sprintf("unsupported record version \"%d\"", RECORD.Version))
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
