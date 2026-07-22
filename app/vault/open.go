package vault

import (
	"os"
	"path/filepath"
	"pvault/app/vault/database"
	v1 "pvault/app/vault/database/database/v1"
	v2 "pvault/app/vault/database/database/v2"
	v3 "pvault/app/vault/database/database/v3"
	"pvault/app/vault/meta"

	"github.com/binarysoupdev/go-commando/errors"
)

func Open(path string) (Vault, error) {
	v, err := detectVault(path)
	if err != nil {
		return Vault{}, err
	}

	v.Map, err = database.LoadIndex(v.Database)
	if err != nil {
		return Vault{}, err
	}

	return v, nil
}

func detectVault(path string) (Vault, error) {
	m, err := detectMetadata(path)
	if err != nil {
		return Vault{}, errors.Chain(err, "error loading vault")
	}

	db, err := loadDatabaseFromVersion(m.DatabaseVersion, path)
	if err != nil {
		return Vault{}, err
	}

	return Vault{
		Meta:     m,
		Database: db,
	}, nil
}

func detectMetadata(path string) (meta.Metadata, error) {
	m := newMetadata(path, 0)

	_, err := os.Stat(m.Path)
	if err == nil {
		return meta.LoadMetadata(m.Path)
	}

	version, ok := detectLegacyDatabaseVersion(path)
	if !ok {
		return meta.Metadata{}, errors.Format("vault not found at \"%s\"", path)
	}

	m.DatabaseVersion = version
	return m, nil
}

func detectLegacyDatabaseVersion(path string) (int, bool) {
	_, err := os.Stat(filepath.Join(path, v2.INDEX_FILE))
	if err == nil {
		return v2.VERSION, true
	}

	_, err = os.Stat(filepath.Join(path, v1.INDEX_FILE))
	if err == nil {
		return v1.VERSION, true
	}

	return 0, false
}

func loadDatabaseFromVersion(version int, path string) (database.Database, error) {
	switch version {
	case v1.VERSION:
		return v1.NewDatabase(path), nil
	case v2.VERSION:
		return v2.NewDatabase(path), nil
	case v3.VERSION:
		return v3.NewDatabase(path), nil
	default:
		return nil, errors.Format("unsupported vault version \"%d\"", version)
	}
}
