package flow

import (
	"pvault/errors"
	"pvault/vault"
)

func OpenVault(path string) (vault.Vault, error) {
	v, err := vault.Open(path)
	if err != nil {
		_, ok := err.(vault.OutOfDateError)
		if ok {
			return vault.Vault{}, errors.New("vault out-of-date (run \"config -upgrade\" to repair)")
		} else {
			return vault.Vault{}, errors.Chain(err, "error opening vault")
		}
	}

	return v, nil
}
