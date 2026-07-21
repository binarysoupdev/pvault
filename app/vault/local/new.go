package local

import (
	"os"
	"pvault/app/vault/database"
	db_v3 "pvault/app/vault/database/version3"
	"pvault/app/vault/index"

	"github.com/binarysoupdev/go-commando/errors"
)

func CreateNewVault(path string) (Vault, error) {
	_, err := os.Stat(path)
	if err == nil || !os.IsNotExist(err) {
		return Vault{}, errors.New("vault path already exists")
	}

	err = os.MkdirAll(path, 0755)
	if err != nil {
		return Vault{}, errors.Chain(err, "error creating vault directory")
	}

	v := Vault{
		Database: db_v3.NewDatabase(path),
		Map:      index.IndexMap{},
	}

	err = v.Database.EncodeIndex(v.Map)
	if err != nil {
		return Vault{}, errors.Chain(err, "error saving initial index")
	}

	return v, nil
}

func OpenVault(path string) (Vault, error) {
	db, err := database.Open(path)
	if err != nil {
		return Vault{}, errors.Chain(err, "error opening database")
	}

	idx, err := db.DecodeIndex()
	if err != nil {
		return Vault{}, errors.Chain(err, "error loading index map")
	}

	return Vault{
		Database: db,
		Map:      idx,
	}, nil
}
