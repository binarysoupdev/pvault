package database

import (
	v2 "pvault/vault/database/version2"
	"pvault/vault/index"

	"github.com/binarysoupdev/go-commando/errors"
)

type Database struct{ Index }

func New(path string) Database {
	return Database{
		Index: v2.NewIndex(path),
	}
}

func Open(path string) (Database, error) {
	idx, err := detectIndex(path)
	if err != nil {
		return Database{}, errors.Chain(err, "error detecting index")
	}

	return Database{
		Index: idx,
	}, nil
}

func (db Database) Initialize(idx index.IndexMap) error {
	err := db.SaveIndex(idx)
	if err != nil {
		return errors.Chain(err, "error saving index file")
	}

	return nil
}
