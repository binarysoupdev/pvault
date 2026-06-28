package flow

import (
	"pvault/errors"
	"pvault/vault"
)

func OpenVault(path string) (vault.Vault, error) {
	v, err := vault.Open(path)
	if err != nil {
		return vault.Vault{}, errors.Chain(err, "error opening vault")
	}

	if v.IsOutOfDate() {
		return vault.Vault{}, errors.Format("vault (@v%d) out-of-date (run \"config -upgrade\" to repair)", v.Version)
	}

	return v, nil
}
