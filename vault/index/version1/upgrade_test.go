package v1_test

import (
	"pvault/vault/data"
	index_v1 "pvault/vault/index/version1"
	"pvault/vault/record"
	record_v1 "pvault/vault/record/version1"
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

	r1 := record_v1.Record{
		Password:      "password1",
		Username:      "username1",
		URL:           "url1",
		RecoveryCodes: []string{"code1-1", "code1-2"},
		ID:            uuid.New(),
		Name:          "name1",
	}
	require.NoError(t, r1.MarshalToLegacy(idx.RecordPath(r1.ID), PASSWORD))

	r2 := record_v1.Record{
		Password:      "password2",
		Username:      "username2",
		URL:           "url2",
		RecoveryCodes: []string{"code2-1", "code2-2"},
		ID:            uuid.New(),
		Name:          "name2",
	}
	require.NoError(t, r2.MarshalToLegacy(idx.RecordPath(r2.ID), PASSWORD))

	MAP := data.NameMap{
		r1.Name: r1.ID,
		r2.Name: r2.ID,
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

	newR1, err := record.Load(newIdx.RecordPath(r1.ID), PASSWORD, r1.ID)
	require.NoError(t, err)
	assert.Equal(t, r1, newR1.(record_v1.Record))

	newR2, err := record.Load(newIdx.RecordPath(r2.ID), PASSWORD, r2.ID)
	require.NoError(t, err)
	assert.Equal(t, r2, newR2.(record_v1.Record))
}
