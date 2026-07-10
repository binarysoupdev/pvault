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
	_, err := os.Stat(filepath.Join(path, version2.INDEX_FILE))
	if err == nil {
		return detectDatabaseFromVersionHeader(path, version2.INDEX_FILE)
	}

	_, err = os.Stat(filepath.Join(path, version1.INDEX_FILE))
	if err == nil {
		return version1.New(path), nil
	}

	return nil, errors.New("index file not found")
}

func detectDatabaseFromVersionHeader(path, indexFile string) (data.Database, error) {
	header := make([]byte, 2)

	file, err := os.Open(filepath.Join(path, indexFile))
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
		return version2.New(path), nil
	default:
		return nil, errors.Format("unsupported version \"%d\"", version)
	}
}
