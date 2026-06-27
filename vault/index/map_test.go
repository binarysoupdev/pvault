package index_test

import (
	"pvault/vault/index"
	"testing"

	"github.com/binarysoupdev/tinsel/rand"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindNameNameNotFoundReturnsNotFound(t *testing.T) {
	//-- arrange
	idx := index.IndexMap{}

	//-- act
	_, found := idx.FindName(uuid.Nil)

	//-- assert
	assert.False(t, found)
}

func TestFindNameNameFoundReturnsName(t *testing.T) {
	//-- arrange
	rand := rand.New(0)
	NAME := rand.ASCII(15)
	ID := uuid.New()

	idx := index.IndexMap{NAME: ID}

	//-- act
	res, found := idx.FindName(ID)

	//-- assert
	assert.True(t, found)
	assert.Equal(t, NAME, res)
}

func TestGetNamesReturnsAllNames(t *testing.T) {
	//-- arrange
	rand := rand.New(0)

	idx := index.IndexMap{
		rand.ASCII(15): uuid.Nil,
		rand.ASCII(15): uuid.Nil,
		rand.ASCII(15): uuid.Nil,
		rand.ASCII(15): uuid.Nil,
		rand.ASCII(15): uuid.Nil,
	}

	//-- act
	res := idx.GetNames()

	//-- assert
	require.Len(t, res, len(idx))
	for name := range idx {
		assert.Contains(t, res, name)
	}
}

func TestSearchNamesReturnsMatches(t *testing.T) {
	//-- arrange
	NAMES := []string{"Foo2", "Bar1", "Foo1"}

	idx := index.IndexMap{
		NAMES[0]: uuid.Nil,
		NAMES[1]: uuid.Nil,
		NAMES[2]: uuid.Nil,
	}

	//-- act
	res := idx.SearchNames("foo")

	//-- assert
	require.Len(t, res, 2)
	assert.Contains(t, res, NAMES[0])
	assert.Contains(t, res, NAMES[2])
}
