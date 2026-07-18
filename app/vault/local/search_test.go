package local_test

import (
	"pvault/app/vault/data"
	vault "pvault/app/vault/local"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSearchEmptyTermReturnsAllNamesSorted(t *testing.T) {
	//-- arrange
	NAMES := []string{"Foo2", "Bar1", "Foo1"}

	v := vault.Vault{
		Index: data.NameMap{
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
	NAMES := []string{"Foo2", "Bar1", "Foo1"}

	v := vault.Vault{
		Index: data.NameMap{
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
