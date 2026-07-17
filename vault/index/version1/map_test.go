package v1_test

import (
	"pvault/vault/data"
	v1 "pvault/vault/index/version1"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveMapLoadMapReturnsMapAndNoError(t *testing.T) {
	//-- arrange
	db := v1.NewIndex(file.NewPath(t, ""))

	MAP := data.NameMap{
		"name1": uuid.New(),
		"name2": uuid.New(),
		"name3": uuid.New(),
	}

	err := db.SaveMap(MAP)
	require.NoError(t, err)

	//-- act
	res, err := db.LoadMap()

	//-- arrange
	require.NoError(t, err)
	assert.Equal(t, MAP, res)
}
