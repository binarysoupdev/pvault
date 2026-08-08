package v1_test

import (
	v1 "pvault/vault/record/record/legacy/v1"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpgradeReturnsUpgradedRecord(t *testing.T) {
	//-- arrange
	r := v1.Record{
		Password:      "password",
		Username:      "username",
		URL:           "url",
		RecoveryCodes: []string{"code1", "code2"},
		ID:            uuid.New(),
		Name:          "name",
	}

	//-- act
	res := r.Upgrade()

	//-- assert
	assert.Equal(t, r.ID, res.ID)
	assert.Equal(t, r.Name, res.Name)
	assert.Equal(t, r.Password, res.Password)
	assert.Equal(t, r.Username, res.Username)
	assert.Equal(t, r.URL, res.Other["url"])
	assert.Equal(t, r.RecoveryCodes, res.Other["recovery_codes"])
}
