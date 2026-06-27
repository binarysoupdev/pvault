package index_test

import (
	"pvault/vault/index"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveLoadIndex(t *testing.T) {
	//-- arrange
	rand := rand.New(0)
	PATH := file.NewPath(t, "")

	idx := index.IndexMap{
		rand.ASCII(10): uuid.New(),
		rand.ASCII(15): uuid.New(),
		rand.ASCII(20): uuid.New(),
	}

	//-- act
	err := idx.Save(PATH)
	require.NoError(t, err)

	res, err := index.LoadIndex(PATH)
	require.NoError(t, err)

	//-- assert
	require.Len(t, res, len(idx))

	for key, val := range idx {
		assert.Contains(t, res, key)
		assert.Equal(t, val, res[key])
	}
}
