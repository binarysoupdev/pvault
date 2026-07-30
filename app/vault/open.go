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
	v := Vault{}
	var err error

	v.Meta, err = loadMetadata(path)
	if err != nil {
		return Vault{}, err
	}

	v.Database, err = loadDatabaseFromVersion(v.Meta.DatabaseVersion, path)
	if err != nil {
		return Vault{}, err
	}

	v.Map, err = database.LoadIndex(v.Database, v.Path)
	if err != nil {
		return Vault{}, err
	}

	return v, nil
}

func loadMetadata(path string) (meta.Metadata, error) {
	metaPath := metadataPath(path)

	_, err := os.Stat(metaPath)
	if err == nil {
		return meta.LoadMetadata(metaPath)
	}

	version, ok := detectLegacyDatabase(path)
	if !ok {
		return meta.Metadata{}, errors.Format("vault not found at \"%s\"", path)
	}
	m := meta.New(version, filepath.Base(path))

	err = meta.SaveMetadata(metaPath, m)
	if err != nil {
		return meta.Metadata{}, errors.Chain(err, "error saving new metadata")
	}

	return m, nil
}

func detectLegacyDatabase(path string) (int, bool) {
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
		return v1.Database{}, nil
	case v2.VERSION:
		return v2.Database{}, nil
	case v3.VERSION:
		return v3.Database{}, nil
	default:
		return nil, errors.Format("unsupported vault version \"%d\"", version)
	}
}
