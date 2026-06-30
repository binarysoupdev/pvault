package legacy

import (
	"pvault/vault/record"

	"github.com/google/uuid"
)

type RecordV1 struct {
	Password      string   `json:"password"`
	Username      string   `json:"username"`
	URL           string   `json:"url"`
	RecoveryCodes []string `json:"recovery_codes"`
}

func (r RecordV1) Upgrade(id uuid.UUID, name string) record.Record {
	other := map[string]any{}

	if r.URL != "" {
		other["url"] = r.URL
	}
	if len(r.RecoveryCodes) > 0 {
		other["recovery_codes"] = r.RecoveryCodes
	}

	return record.Record{
		ID:       id,
		Name:     name,
		Username: r.Username,
		Password: r.Password,
		Other:    other,
	}
}
