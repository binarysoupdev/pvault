package record_test

import (
	"fmt"
	"pvault/vault/record"
	v1 "pvault/vault/record/version1"
	v2 "pvault/vault/record/version2"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadReturnsErrorWhenErrorReadingRecordFile(t *testing.T) {
	//-- arrange
	PATH := file.NewPath(t, "")

	//-- act
	_, res := record.Load(PATH, "", uuid.Nil)

	//-- assert
	require.ErrorContains(t, res, "error reading record file")
}

func TestLoadReturnsErrorWhenVersionUnsupported(t *testing.T) {
	//-- arrange
	VERSION := v2.VERSION + 1

	file, PATH := file.Create(t, "record.json")
	file.Write([]byte{0, byte(VERSION), 0, 0})
	file.Close()

	//-- act
	_, res := record.Load(PATH, "", uuid.Nil)

	//-- assert
	require.ErrorContains(t, res, fmt.Sprintf("unsupported record version \"%d\"", VERSION))
}

func TestLoadReturnsV1RecordAndNoErrorWhenV1Detected(t *testing.T) {
	//-- arrange
	PATH := file.NewPath(t, "record.json")
	const PASSWORD = "Password123!"

	RECORD := v1.Record{}
	require.NoError(t, RECORD.SaveFile(PATH, PASSWORD))

	//-- act
	res, err := record.Load(PATH, PASSWORD, uuid.Nil)

	//-- assert
	require.NoError(t, err)
	assert.IsType(t, RECORD, res)
}

func TestLoadReturnsV2RecordAndNoErrorWhenV2Detected(t *testing.T) {
	//-- arrange
	PATH := file.NewPath(t, "record.json")
	const PASSWORD = "Password123!"

	RECORD := v2.Record{}
	require.NoError(t, RECORD.SaveFile(PATH, PASSWORD))

	//-- act
	res, err := record.Load(PATH, PASSWORD, uuid.Nil)

	//-- assert
	require.NoError(t, err)
	assert.IsType(t, RECORD, res)
}
