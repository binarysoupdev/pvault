package version1_test

import (
	"pvault/vault/data/version1"
	"pvault/vault/index"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadIndexValidLoadsIndexMap(t *testing.T) {
	//-- arrange
	db := version1.New(file.NewPath(t, ""))

	rand := rand.New(0)
	INDEX := index.IndexMap{
		rand.ASCII(10): uuid.New(),
		rand.ASCII(15): uuid.New(),
		rand.ASCII(20): uuid.New(),
	}

	err := db.SaveIndex(INDEX)
	require.NoError(t, err)

	//-- act
	res, err := db.LoadIndex()

	//-- arrange
	require.NoError(t, err)
	assert.Equal(t, INDEX, res)
}
