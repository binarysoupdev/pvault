package v1

import v2 "pvault/vault/record/version2"

func (r Record) Upgrade() v2.Record {
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
