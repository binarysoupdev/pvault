package v1_test

import (
	"os"
	v1 "pvault/vault/database/version/v1"
	v2 "pvault/vault/database/version/v2"
	"pvault/vault/index"
	"pvault/vault/record"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitializeSucceedsAndSavesIndex(t *testing.T) {
	//-- arrange
	db := v1.New(file.NewPath(t, ""))

	INDEX := index.IndexMap{
		"name1": uuid.New(),
		"name2": uuid.New(),
	}

	//-- act
	res := db.Initialize(INDEX)

	//-- assert
	require.NoError(t, res)

	idx, err := db.LoadIndex()
	require.NoError(t, err)
	assert.Equal(t, INDEX, idx)
}

func TestUpgradeValidUpgradesVault(t *testing.T) {
	//-- arrange
	db := v1.New(file.NewPath(t, ""))
	TARGET := v2.New(db.Path)

	file, err := os.Create(db.IndexPath())
	require.NoError(t, err)
	file.Close()

	rand := rand.New(0)
	PASSWORD := rand.ASCII(30)

	LEGACY := map[uuid.UUID]record.RecordV1{}
	INDEX := index.IndexMap{}

	const NUM_LEGACY_FILES = 5
	for range NUM_LEGACY_FILES {
		id := uuid.New()

		INDEX[rand.ASCII(10)] = id
		LEGACY[id] = record.RecordV1{
			Password:      rand.ASCII(30),
			Username:      rand.ASCII(15),
			URL:           rand.ASCII(15),
			RecoveryCodes: []string{rand.ASCII(10), rand.ASCII(10)},
		}

		err := db.SaveLegacyRecord(id, PASSWORD, LEGACY[id])
		require.NoError(t, err)
	}

	//-- act
	res := db.Upgrade(INDEX, TARGET)

	//-- assert
	require.NoError(t, res)
	assert.NoFileExists(t, db.IndexPath())

	for name, id := range INDEX {
		assert.NoFileExists(t, db.RecordPath(id))

		r, err := TARGET.LoadRecord(id, PASSWORD)
		require.NoError(t, err)

		assert.Equal(t, id, r.ID)
		assert.Equal(t, name, r.Name)

		file := LEGACY[id]
		assert.Equal(t, file.Password, r.Password)
		assert.Equal(t, file.Username, r.Username)
		assert.Equal(t, file.URL, r.Other["url"])
		assert.Equal(t, file.RecoveryCodes, r.Other["recovery_codes"])
	}
}
