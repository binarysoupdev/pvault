package v1

import (
	v2 "pvault/vault/record/version2"

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

func (r Record) GetID() uuid.UUID {
	return r.ID
}

func (r Record) Convert() v2.Record {
	other := map[string]any{}

	if r.URL != "" {
		other["url"] = r.URL
	}
	if len(r.RecoveryCodes) > 0 {
		other["recovery_codes"] = r.RecoveryCodes
	}

	return v2.Record{
		ID:       r.ID,
		Name:     r.Name,
		Username: r.Username,
		Password: r.Password,
		Other:    other,
	}
}
