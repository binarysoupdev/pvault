package v1_test

import (
	"pvault/vault/data"
	v1 "pvault/vault/index/version1"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadIndexValidLoadsIndexMap(t *testing.T) {
	//-- arrange
	db := v1.NewIndex(file.NewPath(t, ""))

	rand := rand.New(0)
	INDEX := data.NameMap{
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
