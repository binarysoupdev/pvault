package vault_test

import (
	"pvault/vault"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSearchEmptyTermReturnsAllNamesSorted(t *testing.T) {
	//-- arrange
	NAMES := []string{"foo2", "bar1", "foo1"}

	v := vault.Vault{
		Index: vault.IndexMap{
			NAMES[0]: uuid.New(),
			NAMES[1]: uuid.New(),
			NAMES[2]: uuid.New(),
		},
	}

	//-- act
	res := v.Search("")

	//-- assert
	assert.Equal(t, []string{NAMES[1], NAMES[2], NAMES[0]}, res)
}

func TestSearchReturnsOnlyMatchesSorted(t *testing.T) {
	//-- arrange
	NAMES := []string{"foo2", "bar1", "foo1"}

	v := vault.Vault{
		Index: vault.IndexMap{
			NAMES[0]: uuid.New(),
			NAMES[1]: uuid.New(),
			NAMES[2]: uuid.New(),
		},
	}

	//-- act
	res := v.Search("foo")

	//-- assert
	assert.Equal(t, []string{NAMES[2], NAMES[0]}, res)
}
