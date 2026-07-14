package v1_test

import (
	"os"
	"pvault/vault/data"
	v1 "pvault/vault/index/version1"
	v2 "pvault/vault/index/version2"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpgradeValidUpgradesVault(t *testing.T) {
	//-- arrange
	db := v1.NewIndex(file.NewPath(t, ""))
	TARGET := v2.NewIndex(db.Path)

	file, err := os.Create(db.IndexPath())
	require.NoError(t, err)
	file.Close()

	rand := rand.New(0)
	PASSWORD := rand.ASCII(30)

	LEGACY := map[uuid.UUID]v1.Record{}
	INDEX := data.NameMap{}

	const NUM_LEGACY_FILES = 5
	for range NUM_LEGACY_FILES {
		id := uuid.New()

		INDEX[rand.ASCII(10)] = id
		LEGACY[id] = v1.Record{
			Password:      rand.ASCII(30),
			Username:      rand.ASCII(15),
			URL:           rand.ASCII(15),
			RecoveryCodes: []string{rand.ASCII(10), rand.ASCII(10)},
		}

		file, err := os.Create(db.RecordPath(id))
		require.NoError(t, err)
		defer file.Close()

		LEGACY[id].EncodeToLegacy(file, PASSWORD)
	}

	//-- act
	res := db.Upgrade(INDEX)

	//-- assert
	require.NoError(t, res)
	assert.NoFileExists(t, db.IndexPath())

	for name, id := range INDEX {
		assert.NoFileExists(t, db.RecordPath(id))

		r, err := TARGET.LoadRecord(id, PASSWORD)
		require.NoError(t, err)

		r2 := r.Upgrade()

		assert.Equal(t, id, r2.ID)
		assert.Equal(t, name, r2.Name)

		file := LEGACY[id]
		assert.Equal(t, file.Password, r2.Password)
		assert.Equal(t, file.Username, r2.Username)
		assert.Equal(t, file.URL, r2.Other["url"])
		assert.Equal(t, file.RecoveryCodes, r2.Other["recovery_codes"])
	}
}
