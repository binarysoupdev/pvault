package database

import (
	"os"
	"path/filepath"
	v1 "pvault/app/vault/database/version1"
	v2 "pvault/app/vault/database/version2"
	v3 "pvault/app/vault/database/version3"

	"github.com/binarysoupdev/go-commando/errors"
)

func Open(path string) (Database, error) {
	_, err := os.Stat(filepath.Join(path, v3.INDEX_FILE))
	if err == nil {
		return v3.NewDatabase(path), nil
	}

	_, err = os.Stat(filepath.Join(path, v2.INDEX_FILE))
	if err == nil {
		return v2.NewDatabase(path), nil
	}

	_, err = os.Stat(filepath.Join(path, v1.INDEX_FILE))
	if err == nil {
		return v1.NewDatabase(path), nil
	}

	return nil, errors.Format("index file not found at \"%s\"", path)
}
