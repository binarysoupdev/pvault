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
	v := Vault{
		Path: path,
	}

	err := createOrLoadMetadata(&v)
	if err != nil {
		return Vault{}, err
	}

	v.Database, err = createDatabaseFromVersion(v.Meta.DatabaseVersion)
	if err != nil {
		return Vault{}, err
	}

	err = v.LoadIndex()
	if err != nil {
		return Vault{}, err
	}

	return v, nil
}

func createOrLoadMetadata(v *Vault) error {
	_, err := os.Stat(v.MetadataPath())
	if err == nil {
		return v.LoadMetadata()
	}

	version, ok := detectLegacyDatabase(v.Path)
	if !ok {
		return errors.Format("vault not found at \"%s\"", v.Path)
	}

	v.Meta = meta.New(version, filepath.Base(v.Path))
	return v.SaveMetadata()
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

func createDatabaseFromVersion(version int) (database.Database, error) {
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
