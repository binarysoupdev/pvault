package vault

import (
	"pvault/vault/database"

	"github.com/binarysoupdev/go-commando/errors"
)

func Open(path string) (Vault, error) {
	db, err := database.Find(path)
	if err != nil {
		return Vault{}, err
	}

	idx, err := db.LoadIndex()
	if err != nil {
		return Vault{}, errors.Chain(err, "error parsing index file")
	}

	return Vault{
		Path:     path,
		Index:    idx,
		Database: db,
	}, nil
}
