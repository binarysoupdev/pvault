package v1_test

import (
	record "pvault/vault/record/version1"
	"testing"

	"github.com/binarysoupdev/tinsel/rand"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRecordConvertConvertsRecord(t *testing.T) {
	//-- arrange
	rand := rand.New(0)

	r := record.Record{
		Password:      rand.ASCII(30),
		Username:      rand.ASCII(15),
		URL:           rand.ASCII(15),
		RecoveryCodes: []string{rand.ASCII(10), rand.ASCII(10)},
		ID:            uuid.New(),
		Name:          rand.ASCII(15),
	}

	//-- act
	res := r.Convert()

	//-- assert
	assert.Equal(t, r.ID, res.ID)
	assert.Equal(t, r.Name, res.Name)
	assert.Equal(t, r.Password, res.Password)
	assert.Equal(t, r.Username, res.Username)
	assert.Equal(t, r.URL, res.Other["url"])
	assert.Equal(t, r.RecoveryCodes, res.Other["recovery_codes"])
}
