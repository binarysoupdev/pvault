package v1_test

import (
	v1 "pvault/vault/record/version/v1"
	"testing"

	"github.com/binarysoupdev/tinsel/rand"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRecordUpgradeUpgradesRecord(t *testing.T) {
	//-- arrange
	rand := rand.New(0)

	ID := uuid.New()
	NAME := rand.ASCII(15)

	r := v1.Record{
		Password:      rand.ASCII(30),
		Username:      rand.ASCII(15),
		URL:           rand.ASCII(15),
		RecoveryCodes: []string{rand.ASCII(10), rand.ASCII(10)},
	}

	//-- act
	res := r.Upgrade(ID, NAME)

	//-- assert
	assert.Equal(t, ID, res.ID)
	assert.Equal(t, NAME, res.Name)
	assert.Equal(t, r.Password, res.Password)
	assert.Equal(t, r.Username, res.Username)
	assert.Equal(t, r.URL, res.Other["url"])
	assert.Equal(t, r.RecoveryCodes, res.Other["recovery_codes"])
}
