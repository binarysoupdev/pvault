package v2_test

import (
	record "pvault/vault/record/version2"
	"testing"

	"github.com/binarysoupdev/tinsel/rand"
	"github.com/stretchr/testify/assert"
)

func TestValidateInvalid(t *testing.T) {
	//-- act
	res := record.Record{}.Validate()

	//-- assert
	assert.ErrorContains(t, res, "\"ID\" cannot be nil (all zeroes)")
	assert.ErrorContains(t, res, "\"Name\" cannot be empty")
}

func TestValidateValid(t *testing.T) {
	//-- arrange
	rand := rand.New(0)
	r := record.NewEmptyRecord(rand.ASCII(15))

	//-- act
	res := r.Validate()

	//-- assert
	assert.NoError(t, res)
}
