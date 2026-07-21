package v1_test

import (
	"pvault/app/vault/index"
	index_v1 "pvault/app/vault/index/version1"
	"pvault/app/vault/record"
	record_v1 "pvault/app/vault/record/version1"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpgradeReturnsNewIndexAndNoErrorAndUpgradesVault(t *testing.T) {
	//-- arrange
	idx := index_v1.NewIndex(file.NewPath(t, ""))
	const PASSWORD = "Password123!"

	R1 := record_v1.Record{
		Password:      "password1",
		Username:      "username1",
		URL:           "url1",
		RecoveryCodes: []string{"code1-1", "code1-2"},
		ID:            uuid.New(),
		Name:          "name1",
	}
	require.NoError(t, R1.SaveLegacy(idx.RecordPath(R1.ID), PASSWORD))

	R2 := record_v1.Record{
		Password:      "password2",
		Username:      "username2",
		URL:           "url2",
		RecoveryCodes: []string{"code2-1", "code2-2"},
		ID:            uuid.New(),
		Name:          "name2",
	}
	require.NoError(t, R2.SaveLegacy(idx.RecordPath(R2.ID), PASSWORD))

	MAP := index.IndexMap{
		R1.Name: R1.ID,
		R2.Name: R2.ID,
	}
	require.NoError(t, idx.SaveMap(MAP))

	//-- act
	newIdx, err := idx.Upgrade()

	//-- assert
	require.NoError(t, err)
	assert.NoFileExists(t, idx.Filepath())

	m, err := newIdx.LoadMap()
	require.NoError(t, err)
	assert.Equal(t, m, MAP)

	r1, err := record.Load(newIdx.RecordPath(R1.ID), PASSWORD, R1.ID)
	require.NoError(t, err)
	assert.Equal(t, R1, r1.(record_v1.Record))

	r2, err := record.Load(newIdx.RecordPath(R2.ID), PASSWORD, R2.ID)
	require.NoError(t, err)
	assert.Equal(t, R2, r2.(record_v1.Record))
}
