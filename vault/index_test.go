package vault_test

import (
	"pvault/vault"
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
	PATH := file.NewPath(t, rand.ASCII(8))

	INDEX := vault.IndexMap{
		rand.ASCII(10): uuid.New(),
		rand.ASCII(15): uuid.New(),
		rand.ASCII(20): uuid.New(),
	}

	//-- act
	err := INDEX.Save(PATH)
	require.NoError(t, err)

	res, err := vault.LoadIndex(PATH)
	require.NoError(t, err)

	//-- assert
	require.Len(t, res, len(INDEX))

	for key, val := range INDEX {
		assert.Contains(t, res, key)
		assert.Equal(t, val, res[key])
	}
}
