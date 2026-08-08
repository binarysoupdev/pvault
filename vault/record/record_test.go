package record_test

import (
	"pvault/vault/record"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateReturnsErrorWhenInvalid(t *testing.T) {
	//-- arrange
	r := &record.Mock{
		ID:   uuid.Nil,
		Name: "",
	}

	//-- act
	res := record.Validate(r)

	//-- assert
	assert.ErrorContains(t, res, "id cannot be nil (all zeroes)")
	assert.ErrorContains(t, res, "name cannot be empty")
}

func TestValidateReturnsNoErrorWhenValid(t *testing.T) {
	//-- arrange
	r := &record.Mock{
		ID:   uuid.New(),
		Name: "name",
	}

	//-- act
	res := record.Validate(r)

	//-- assert
	require.NoError(t, res)
}
