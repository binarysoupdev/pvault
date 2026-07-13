package record_test

import (
	"pvault/vault/record"
	"testing"

	"github.com/binarysoupdev/tinsel/rand"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRecordV1UpgradeUpgradesRecord(t *testing.T) {
	//-- arrange
	rand := rand.New(0)

	ID := uuid.New()
	NAME := rand.ASCII(15)

	v1 := record.RecordV1{
		Password:      rand.ASCII(30),
		Username:      rand.ASCII(15),
		URL:           rand.ASCII(15),
		RecoveryCodes: []string{rand.ASCII(10), rand.ASCII(10)},
	}

	//-- act
	res := v1.Upgrade(ID, NAME)

	//-- assert
	assert.Equal(t, ID, res.ID)
	assert.Equal(t, NAME, res.Name)
	assert.Equal(t, v1.Password, res.Password)
	assert.Equal(t, v1.Username, res.Username)
	assert.Equal(t, v1.URL, res.Other["url"])
	assert.Equal(t, v1.RecoveryCodes, res.Other["recovery_codes"])
}
