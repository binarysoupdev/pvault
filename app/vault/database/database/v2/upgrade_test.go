package v2_test

import (
	"encoding/binary"
	"fmt"
	"os"
	"pvault/app/vault/database"
	db_v2 "pvault/app/vault/database/database/v2"
	"pvault/app/vault/index"
	record_v1 "pvault/app/vault/record/record/v1"
	record_v2 "pvault/app/vault/record/record/v2"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpgradeReturnsErrorWhenOldIndexNotFound(t *testing.T) {
	//-- act
	_, res := db_v2.Database{}.Upgrade("invalid")

	//-- assert
	assert.ErrorContains(t, res, "error loading old index file")
}

func TestUpgradeReturnsErrorWhenErrorBackingRecords(t *testing.T) {
	//-- arrange
	db := db_v2.Database{}
	PATH := file.NewPath(t, "")

	R1 := uuid.New()
	R2 := uuid.New()

	INDEX := index.IndexMap{
		"name1": R1,
		"name2": R2,
	}
	require.NoError(t, database.SaveIndex(db, PATH, INDEX))

	//-- act
	_, res := db_v2.Database{}.Upgrade(PATH)

	//-- assert
	assert.ErrorContains(t, res, "error backing record "+R1.String())
	assert.ErrorContains(t, res, "error backing record "+R2.String())
}

func TestUpgradeReturnsErrorWhenUnsupportedRecordVersionDetected(t *testing.T) {
	//-- arrange
	db := db_v2.Database{}
	PATH := file.NewPath(t, "")

	const VERSION = 3
	ID := uuid.New()

	version := make([]byte, 2)
	binary.BigEndian.PutUint16(version, VERSION)
	require.NoError(t, os.WriteFile(db.RecordPath(PATH, ID), version, 0666))

	require.NoError(t, database.SaveIndex(db, PATH, index.IndexMap{"": ID}))

	//-- act
	_, res := db.Upgrade(PATH)

	//-- assert
	assert.ErrorContains(t, res, fmt.Sprintf("unsupported record version \"%d\"", VERSION))
}

func TestUpgradeReturnsNewDatabaseAndNoErrorAndUpgradesDatabase(t *testing.T) {
	//-- arrange
	db := db_v2.Database{}
	PATH := file.NewPath(t, "")
	const PASSWORD = "Password123!"

	R1 := record_v1.Record{
		ID:            uuid.New(),
		Name:          "name1",
		Password:      "password1",
		Username:      "username1",
		URL:           "url1",
		RecoveryCodes: []string{"code1-1", "code1-2"},
	}
	require.NoError(t, database.SaveRecord(db, PATH, R1, PASSWORD))

	R2 := record_v2.Record{
		ID:       uuid.New(),
		Name:     "name2",
		Password: "password2",
		Username: "username2",
	}
	require.NoError(t, database.SaveRecord(db, PATH, R2, PASSWORD))

	INDEX := index.IndexMap{
		R1.Name: R1.ID,
		R2.Name: R2.ID,
	}
	require.NoError(t, database.SaveIndex(db, PATH, INDEX))

	//-- act
	res, err := db.Upgrade(PATH)

	//-- assert
	require.NoError(t, err)
	assert.NoFileExists(t, db.IndexPath(PATH))

	idx, err := database.LoadIndex(res, PATH)
	require.NoError(t, err)
	assert.Equal(t, idx, INDEX)

	r1, err := database.LoadRecord(res, PATH, R1.ID, PASSWORD)
	require.NoError(t, err)
	assert.Equal(t, R1, r1)

	r2, err := database.LoadRecord(res, PATH, R2.ID, PASSWORD)
	require.NoError(t, err)
	assert.Equal(t, R2, r2)
}
