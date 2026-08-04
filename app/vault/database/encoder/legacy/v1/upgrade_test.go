package v1_test

import (
	"pvault/app/vault/database"
	db_v1 "pvault/app/vault/database/encoder/legacy/v1"
	"pvault/app/vault/index"
	record_v1 "pvault/app/vault/record/record/legacy/v1"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpgradeReturnsErrorWhenOldIndexNotFound(t *testing.T) {
	//-- act
	_, res := db_v1.Encoder{}.Upgrade("invalid")

	//-- assert
	assert.ErrorContains(t, res, "error loading old index file")
}

func TestUpgradeReturnsErrorWhenErrorBackingRecords(t *testing.T) {
	//-- arrange
	db := db_v1.Encoder{}
	PATH := t.TempDir()

	R1 := uuid.New()
	R2 := uuid.New()

	INDEX := index.IndexMap{
		"name1": R1,
		"name2": R2,
	}
	require.NoError(t, database.SaveIndex(db, PATH, INDEX))

	//-- act
	_, res := db_v1.Encoder{}.Upgrade(PATH)

	//-- assert
	assert.ErrorContains(t, res, "error backing record "+R1.String())
	assert.ErrorContains(t, res, "error backing record "+R2.String())
}

func TestUpgradeReturnsNewEncoderAndNoErrorAndUpgradesEncoder(t *testing.T) {
	//-- arrange
	db := db_v1.Encoder{}
	PATH := t.TempDir()
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

	R2 := record_v1.Record{
		ID:            uuid.New(),
		Name:          "name2",
		Password:      "password2",
		Username:      "username2",
		URL:           "url2",
		RecoveryCodes: []string{"code2-1", "code2-2"},
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
	assert.NoFileExists(t, db.RecordPath(PATH, R1.ID))

	r2, err := database.LoadRecord(res, PATH, R2.ID, PASSWORD)
	require.NoError(t, err)
	assert.Equal(t, R2, r2)
	assert.NoFileExists(t, db.RecordPath(PATH, R2.ID))
}
