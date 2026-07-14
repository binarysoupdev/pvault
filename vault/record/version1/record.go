package v1

import (
	"github.com/google/uuid"
)

const VERSION = 1

type Record struct {
	Password      string   `json:"password"`
	Username      string   `json:"username"`
	URL           string   `json:"url"`
	RecoveryCodes []string `json:"recovery_codes"`

	ID   uuid.UUID
	Name string
}

func (r Record) GetVersion() int {
	return VERSION
}

func (r Record) GetID() uuid.UUID {
	return r.ID
}

func (r Record) GetName() string {
	return r.Name
}
