package v2_test

import (
	v2 "pvault/vault/record/version/v2"
	"testing"

	"github.com/binarysoupdev/tinsel/rand"
	"github.com/stretchr/testify/assert"
)

func TestValidateInvalid(t *testing.T) {
	//-- act
	res := v2.Record{}.Validate()

	//-- assert
	assert.ErrorContains(t, res, "\"ID\" cannot be nil (all zeroes)")
	assert.ErrorContains(t, res, "\"Name\" cannot be empty")
}

func TestValidateValid(t *testing.T) {
	//-- arrange
	rand := rand.New(0)
	r := v2.NewFromName(rand.ASCII(15))

	//-- act
	res := r.Validate()

	//-- assert
	assert.NoError(t, res)
}
