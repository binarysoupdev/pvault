package index_test

import (
	"pvault/vault/index"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindNameReturnsNotFoundWhenNameNotFound(t *testing.T) {
	//-- arrange
	m := index.IndexMap{}

	//-- act
	_, found := m.FindName(uuid.Nil)

	//-- assert
	assert.False(t, found)
}

func TestFindNameReturnsNameWhenNameFound(t *testing.T) {
	//-- arrange
	NAME := "name"
	ID := uuid.New()

	m := index.IndexMap{NAME: ID}

	//-- act
	res, found := m.FindName(ID)

	//-- assert
	assert.True(t, found)
	assert.Equal(t, NAME, res)
}

func TestGetNamesReturnsAllNames(t *testing.T) {
	//-- arrange
	m := index.IndexMap{
		"name1": uuid.Nil,
		"name2": uuid.Nil,
		"name3": uuid.Nil,
	}

	//-- act
	res := m.GetNames()

	//-- assert
	require.Len(t, res, len(m))
	for name := range m {
		assert.Contains(t, res, name)
	}
}

func TestSearchNamesReturnsMatches(t *testing.T) {
	//-- arrange
	NAMES := []string{"Foo2", "Bar1", "Foo1"}

	m := index.IndexMap{
		NAMES[0]: uuid.Nil,
		NAMES[1]: uuid.Nil,
		NAMES[2]: uuid.Nil,
	}

	//-- act
	res := m.SearchNames("foo")

	//-- assert
	require.Len(t, res, 2)
	assert.Contains(t, res, NAMES[0])
	assert.Contains(t, res, NAMES[2])
}
