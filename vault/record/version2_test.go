package record_test

import (
	"pvault/vault/record"
	"testing"

	"github.com/binarysoupdev/tinsel/rand"
	"github.com/stretchr/testify/assert"
)

func TestValidateInvalid(t *testing.T) {
	//-- act
	res := record.RecordV2{}.Validate()

	//-- assert
	assert.ErrorContains(t, res, "\"ID\" cannot be nil (all zeroes)")
	assert.ErrorContains(t, res, "\"Name\" cannot be empty")
}

func TestValidateValid(t *testing.T) {
	//-- arrange
	rand := rand.New(0)
	r := record.NewFromName(rand.ASCII(15))

	//-- act
	res := r.Validate()

	//-- assert
	assert.NoError(t, res)
}
