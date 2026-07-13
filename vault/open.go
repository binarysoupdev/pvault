package vault

import (
	"pvault/vault/database/version"

	"github.com/binarysoupdev/go-commando/errors"
)

func Open(path string) (Vault, error) {
	database, err := version.Detect(path)
	if err != nil {
		return Vault{}, err
	}

	idx, err := database.LoadIndex()
	if err != nil {
		return Vault{}, errors.Chain(err, "error parsing index file")
	}

	return Vault{
		Path:     path,
		Index:    idx,
		Database: database,
	}, nil
}
