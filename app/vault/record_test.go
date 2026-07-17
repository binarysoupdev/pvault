package vault_test

import (
	"fmt"
	"pvault/app/vault"
	"pvault/app/vault/data"
	record "pvault/app/vault/record/version2"
	"testing"

	"github.com/binarysoupdev/tinsel/rand"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidateRecordInvalidRecord(t *testing.T) {
	//-- arrange
	RECORD := record.NewEmptyRecord("")

	//-- act
	res := vault.Vault{}.ValidateRecord(RECORD)

	//-- assert
	assert.ErrorContains(t, res, "\"Name\" cannot be empty")
}

func TestValidateRecordNameAlreadyExists(t *testing.T) {
	//-- arrange
	rand := rand.New(0)
	NAME := rand.ASCII(15)

	RECORD := record.NewEmptyRecord(NAME)

	v := vault.Vault{
		Index: data.NameMap{
			NAME: uuid.Nil,
		},
	}

	//-- act
	res := v.ValidateRecord(RECORD)

	//-- assert
	assert.ErrorContains(t, res, fmt.Sprintf("name \"%s\" already exists", NAME))
}
