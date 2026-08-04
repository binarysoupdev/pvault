package vault_test

import (
	"fmt"
	"os"
	"pvault/app/vault"
	"pvault/app/vault/database"
	"pvault/app/vault/index"
	"pvault/app/vault/record"
	record_v1 "pvault/app/vault/record/record/legacy/v1"
	record_v2 "pvault/app/vault/record/record/v2"
	"testing"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateRecordReturnsErrorWhenRecordInvalid(t *testing.T) {
	//-- act
	res := vault.Vault{}.ValidateRecord(&record.Mock{})

	//-- assert
	assert.ErrorContains(t, res, "id cannot be nil")
}

func TestValidateRecordReturnsErrorWhenNameAlreadyExistsForAnotherRecord(t *testing.T) {
	//-- arrange
	RECORD := record.NewMock("name")

	v := vault.Vault{
		Map: index.IndexMap{
			RECORD.Name: uuid.New(),
		},
	}

	//-- act
	res := v.ValidateRecord(RECORD)

	//-- assert
	assert.ErrorContains(t, res, fmt.Sprintf("name \"%s\" already exists", RECORD.Name))
}

func TestValidateRecordReturnsNoErrorWhenNameExistsForSameRecord(t *testing.T) {
	//-- arrange
	RECORD := record.NewMock("name")

	v := vault.Vault{
		Map: index.IndexMap{
			RECORD.Name: RECORD.ID,
		},
	}

	//-- act
	res := v.ValidateRecord(RECORD)

	//-- assert
	require.NoError(t, res)
}

func TestSaveRecordReturnsErrorWhenRecordInvalid(t *testing.T) {
	//-- act
	res := vault.Vault{}.SaveRecord(&record.Mock{}, "")

	//-- assert
	assert.ErrorContains(t, res, "error validating record")
}

func TestSaveRecordReturnsErrorWhenDatabaseSaveRecordReturnsError(t *testing.T) {
	//-- arrange
	v := vault.Vault{
		Path: t.TempDir(),
		Database: &database.Mock{
			EncodeRecordError: errors.New(""),
		},
	}

	//-- act
	res := v.SaveRecord(record.NewMock("name"), "")

	//-- assert
	assert.ErrorContains(t, res, "error saving record")
}

func TestSaveRecordReturnsErrorWhenSaveIndexReturnsError(t *testing.T) {
	//-- arrange
	v := vault.Vault{
		Path: t.TempDir(),
		Map:  index.IndexMap{},
		Database: &database.Mock{
			EncodeIndexError: errors.New(""),
		},
	}

	//-- act
	res := v.SaveRecord(record.NewMock("name"), "")

	//-- assert
	assert.ErrorContains(t, res, "error saving index map")
}

func TestSaveRecordReturnsNoErrorWithNewIdAndNewName(t *testing.T) {
	//-- arrange
	RECORD := record.NewMock("name")
	mock := &database.Mock{}

	v := vault.Vault{
		Path:     t.TempDir(),
		Database: mock,
		Map:      index.IndexMap{},
	}

	//-- act
	res := v.SaveRecord(RECORD, "")

	//-- assert
	require.NoError(t, res)
	assert.Equal(t, RECORD, mock.Record)

	assert.Contains(t, v.Map, RECORD.Name)
}

func TestSaveRecordReturnsNoErrorWithExistingIDAndNewName(t *testing.T) {
	//-- arrange
	RECORD := record.NewMock("name")
	mock := &database.Mock{}

	v := vault.Vault{
		Path:     t.TempDir(),
		Database: mock,
		Map: index.IndexMap{
			RECORD.Name: RECORD.ID,
		},
	}
	RECORD.Name += "x"

	//-- act
	res := v.SaveRecord(RECORD, "")

	//-- assert
	require.NoError(t, res)
	assert.Equal(t, RECORD, mock.Record)

	assert.Contains(t, v.Map, RECORD.Name)
	assert.Len(t, v.Map, 1)
}

func TestSaveRecordReturnsErrorWithNewIDAndExistingName(t *testing.T) {
	//-- arrange
	RECORD := record.NewMock("name")

	v := vault.Vault{
		Path:     t.TempDir(),
		Database: &database.Mock{},
		Map: index.IndexMap{
			RECORD.Name: uuid.New(),
		},
	}

	//-- act
	res := v.SaveRecord(RECORD, "")

	//-- assert
	assert.ErrorContains(t, res, "error validating record")
}

func TestSaveRecordReturnsNoErrorWithExistingIDAndExistingName(t *testing.T) {
	//-- arrange
	RECORD := record.NewMock("name")
	mock := &database.Mock{}

	v := vault.Vault{
		Path:     t.TempDir(),
		Database: mock,
		Map: index.IndexMap{
			RECORD.Name: RECORD.ID,
		},
	}

	//-- act
	res := v.SaveRecord(RECORD, "")

	//-- assert
	require.NoError(t, res)
	assert.Equal(t, RECORD, mock.Record)

	assert.Contains(t, v.Map, RECORD.Name)
	assert.Len(t, v.Map, 1)
}

func TestLoadRecordReturnsErrorWhenNameNotFound(t *testing.T) {
	//-- arrange
	RECORD := record.NewMock("name")

	//-- act
	_, res := vault.Vault{}.LoadRecord(RECORD.Name, "")

	//-- assert
	assert.ErrorContains(t, res, fmt.Sprintf("name \"%s\" not found", RECORD.Name))
}

func TestLoadRecordReturnsErrorWhenDatabaseLoadRecordReturnsError(t *testing.T) {
	//-- arrange
	RECORD := record.NewMock("name")
	mock := &database.Mock{
		DecodeRecordError: errors.New(""),
	}

	v := vault.Vault{
		Path:     t.TempDir(),
		Database: mock,
		Map: index.IndexMap{
			RECORD.Name: RECORD.ID,
		},
	}
	require.NoError(t, os.WriteFile(mock.RecordPath(v.Path, RECORD.ID), []byte{}, 0666))

	//-- act
	_, res := v.LoadRecord(RECORD.Name, "")

	//-- assert
	assert.ErrorContains(t, res, "error loading record")
}

func TestLoadRecordReturnsV1RecordAndNoError(t *testing.T) {
	//-- arrange
	RECORD := record_v1.Record{
		ID:   uuid.New(),
		Name: "name",
	}
	mock := &database.Mock{
		Record: RECORD,
	}

	v := vault.Vault{
		Path:     t.TempDir(),
		Database: mock,
		Map: index.IndexMap{
			RECORD.Name: RECORD.ID,
		},
	}
	require.NoError(t, os.WriteFile(mock.RecordPath(v.Path, RECORD.ID), []byte{}, 0666))

	//-- act
	res, err := v.LoadRecord(RECORD.Name, "")

	//-- assert
	require.NoError(t, err)
	assert.Equal(t, RECORD, res)
}

func TestLoadRecordReturnsV2RecordAndNoError(t *testing.T) {
	//-- arrange
	RECORD := record_v2.Record{
		ID:   uuid.New(),
		Name: "name",
	}
	mock := &database.Mock{
		Record: RECORD,
	}

	v := vault.Vault{
		Path:     t.TempDir(),
		Database: mock,
		Map: index.IndexMap{
			RECORD.Name: RECORD.ID,
		},
	}
	require.NoError(t, os.WriteFile(mock.RecordPath(v.Path, RECORD.ID), []byte{}, 0666))

	//-- act
	res, err := v.LoadRecord(RECORD.Name, "")

	//-- assert
	require.NoError(t, err)
	assert.Equal(t, RECORD, res)
}

func TestDeleteRecordReturnsErrorWhenNameNotFound(t *testing.T) {
	//-- arrange
	RECORD := record.NewMock("name")

	//-- act
	_, res := vault.Vault{}.DeleteRecord(RECORD.Name)

	//-- assert
	assert.ErrorContains(t, res, fmt.Sprintf("name \"%s\" not found", RECORD.Name))
}

func TestDeleteRecordReturnsErrorWhenDatabaseDeleteRecordReturnsError(t *testing.T) {
	//-- arrange
	RECORD := record.NewMock("name")

	v := vault.Vault{
		Path:     t.TempDir(),
		Database: &database.Mock{},
		Map: index.IndexMap{
			RECORD.Name: RECORD.ID,
		},
	}

	//-- act
	_, res := v.DeleteRecord(RECORD.Name)

	//-- assert
	assert.ErrorContains(t, res, "error deleting record")
}

func TestDeleteRecordReturnsErrorWhenSaveIndexReturnsError(t *testing.T) {
	//-- arrange
	RECORD := record.NewMock("name")
	mock := &database.Mock{
		EncodeIndexError: errors.New(""),
	}

	v := vault.Vault{
		Path:     t.TempDir(),
		Database: mock,
		Map: index.IndexMap{
			RECORD.Name: RECORD.ID,
		},
	}
	require.NoError(t, os.WriteFile(mock.RecordPath(v.Path, RECORD.ID), []byte{}, 0666))

	//-- act
	_, res := v.DeleteRecord(RECORD.Name)

	//-- assert
	assert.ErrorContains(t, res, "error saving index map")
}

func TestDeleteRecordReturnsIDAndNoErrorAndDeletesRecord(t *testing.T) {
	//-- arrange
	RECORD := record.NewMock("name")
	mock := &database.Mock{}

	v := vault.Vault{
		Path:     t.TempDir(),
		Database: mock,
		Map: index.IndexMap{
			RECORD.Name: RECORD.ID,
		},
	}
	require.NoError(t, os.WriteFile(mock.RecordPath(v.Path, RECORD.ID), []byte{}, 0666))

	//-- act
	id, err := v.DeleteRecord(RECORD.Name)

	//-- assert
	require.NoError(t, err)
	assert.Equal(t, RECORD.ID, id)

	assert.Empty(t, v.Map)
	assert.NoFileExists(t, mock.RecordPath(v.Path, RECORD.ID))
}
