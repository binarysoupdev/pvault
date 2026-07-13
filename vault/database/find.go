package database

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"pvault/vault/database/version1"
	"pvault/vault/database/version2"

	"github.com/binarysoupdev/go-commando/errors"
)

func Find(path string) (Database, error) {
	_, err := os.Stat(filepath.Join(path, version2.INDEX_FILE))
	if err == nil {
		return detectFromVersionHeader(path, version2.INDEX_FILE)
	}

	_, err = os.Stat(filepath.Join(path, version1.INDEX_FILE))
	if err == nil {
		return version1.NewDatabase(path), nil
	}

	return nil, errors.New("index file not found")
}

func detectFromVersionHeader(path, indexFile string) (Database, error) {
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
		return version2.NewDatabase(path), nil
	default:
		return nil, errors.Format("unsupported version \"%d\"", version)
	}
}
