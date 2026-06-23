package vault_test

import (
	"fmt"
	"pvault/vault"
	"pvault/vault/index"
	"pvault/vault/record"
	"testing"

	"github.com/binarysoupdev/tinsel/rand"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidateRecordInvalidRecord(t *testing.T) {
	//-- arrange
	RECORD := record.NewFromName("")

	//-- act
	res := vault.Vault{}.ValidateRecord(RECORD)

	//-- assert
	assert.ErrorContains(t, res, "\"Name\" cannot be empty")
}

func TestValidateRecordNameAlreadyExists(t *testing.T) {
	//-- arrange
	rand := rand.New(0)
	NAME := rand.ASCII(15)

	RECORD := record.NewFromName(NAME)

	v := vault.Vault{
		Index: index.IndexMap{
			NAME: uuid.Nil,
		},
	}

	//-- act
	res := v.ValidateRecord(RECORD)

	//-- assert
	assert.ErrorContains(t, res, fmt.Sprintf("name \"%s\" already exists", NAME))
}
