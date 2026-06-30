package vault

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"pvault/errors"
	"pvault/vault/data"
	"pvault/vault/data/version1"
	"pvault/vault/data/version2"
)

const LEGACY_INDEX_FILE = "index.txt"

func Open(path string) (Vault, error) {
	db, err := detectDatabase(path)
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

func detectDatabase(path string) (data.Database, error) {
	indexPath := filepath.Join(path, INDEX_FILE)

	_, err := os.Stat(indexPath)
	if err == nil {
		return detectDatabaseFromVersionHeader(indexPath)
	}

	legacyPath := filepath.Join(path, LEGACY_INDEX_FILE)

	_, err = os.Stat(legacyPath)
	if err == nil {
		return version1.NewDatabase(legacyPath), nil
	}

	return nil, errors.New("index file not found")
}

func detectDatabaseFromVersionHeader(path string) (data.Database, error) {
	header := make([]byte, 2)

	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Chain(err, "error opening index file")
	}
	defer file.Close()

	_, err = file.Read(header)
	if err != nil {
		return nil, errors.Chain(err, "error reading version header")
	}

	version := binary.BigEndian.Uint16(header)

	switch version {
	case 2:
		return version2.NewDatabase(path), nil
	default:
		return nil, errors.Format("unsupported version \"%d\"", version)
	}
}
