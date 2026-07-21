package vault

import (
	"os"
	"pvault/app/vault/database"
	db_v3 "pvault/app/vault/database/version3"
	"pvault/app/vault/index"

	"github.com/binarysoupdev/go-commando/errors"
)

func CreateNew(path string) (Vault, error) {
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

	err = database.SaveIndex(v.Database, v.Map)
	if err != nil {
		return Vault{}, errors.Chain(err, "error saving initial index")
	}

	return v, nil
}

func Open(path string) (Vault, error) {
	db, err := database.Open(path)
	if err != nil {
		return Vault{}, errors.Chain(err, "error opening database")
	}

	idx, err := database.LoadIndex(db)
	if err != nil {
		return Vault{}, errors.Chain(err, "error loading index")
	}

	return Vault{
		Database: db,
		Map:      idx,
	}, nil
}
