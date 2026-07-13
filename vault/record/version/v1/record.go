package v1

import (
	v2 "pvault/vault/record/version/v2"

	"github.com/google/uuid"
)

const VERSION = 1

type Record struct {
	Password      string   `json:"password"`
	Username      string   `json:"username"`
	URL           string   `json:"url"`
	RecoveryCodes []string `json:"recovery_codes"`
}

func (r Record) Upgrade(id uuid.UUID, name string) v2.Record {
	other := map[string]any{}

	if r.URL != "" {
		other["url"] = r.URL
	}
	if len(r.RecoveryCodes) > 0 {
		other["recovery_codes"] = r.RecoveryCodes
	}

	return v2.Record{
		ID:       id,
		Name:     name,
		Username: r.Username,
		Password: r.Password,
		Other:    other,
	}
}

func Downgrade(r v2.Record) Record {
	return Record{
		Password: r.Password,
		Username: r.Username,
	}
}
