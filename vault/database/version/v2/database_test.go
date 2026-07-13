package v2_test

import (
	v2 "pvault/vault/database/version/v2"
	"pvault/vault/index"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitializeSucceedsAndSavesIndex(t *testing.T) {
	//-- arrange
	db := v2.New(file.NewPath(t, ""))

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
