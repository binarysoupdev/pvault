package database

import (
	"os"
	"pvault/app/vault/index"

	"github.com/binarysoupdev/go-commando/errors"
)

func SaveIndex(db Database, idx index.IndexMap) error {
	file, err := os.Create(db.IndexPath())
	if err != nil {
		return errors.Chain(err, "error creating index file")
	}
	defer file.Close()

	err = db.EncodeIndex(file, idx)
	if err != nil {
		return errors.Chain(err, "error encoding index")
	}

	return nil
}

func LoadIndex(db Database) (index.IndexMap, error) {
	file, err := os.Open(db.IndexPath())
	if err != nil {
		return nil, errors.Chain(err, "error opening index file")
	}
	defer file.Close()

	idx, err := db.DecodeIndex(file)
	if err != nil {
		return nil, errors.Chain(err, "error decoding index")
	}

	return idx, nil
}
