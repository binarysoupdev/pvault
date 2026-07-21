package database

import (
	"pvault/app/vault/index"
	"pvault/app/vault/index/encoder"
)

func SaveIndex(db Database, idx index.IndexMap) error {
	return encoder.SaveIndexFile(db, idx, db.IndexPath())
}

func LoadIndex(db Database) (index.IndexMap, error) {
	return encoder.LoadIndexFile(db, db.IndexPath())
}
