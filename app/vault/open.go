package vault

import (
	"os"
	"path/filepath"
	"pvault/app/vault/database"
	db_v1 "pvault/app/vault/database/database/legacy/v1"
	db_v2 "pvault/app/vault/database/database/legacy/v2"
	db_v3 "pvault/app/vault/database/database/v3"
	"pvault/app/vault/meta"
	meta_encoder "pvault/app/vault/meta/encoder"

	"github.com/binarysoupdev/go-commando/errors"
)

func Open(path string) (Vault, error) {
	v := New(path, filepath.Base(path))

	err := open(&v)
	if err != nil {
		return Vault{}, err
	}

	err = v.LoadIndex()
	if err != nil {
		return Vault{}, err
	}

	return v, nil
}

func open(v *Vault) error {
	ok, err := loadFromMetadata(v)
	if ok || err != nil {
		return err
	}

	ok, err = loadFromDatabase(v)
	if ok || err != nil {
		return err
	}

	return errors.Format("vault not found at \"%s\"", v.Path)
}

func loadFromMetadata(v *Vault) (bool, error) {
	encoder := detectMetadata(v.Path)
	if encoder == nil {
		return false, nil
	}

	v.MetaEncoder = encoder

	err := v.LoadMetadata()
	if err != nil {
		return false, err
	}

	switch v.Meta.DatabaseVersion {
	case db_v1.VERSION:
		v.Database = db_v1.Database{}
	case db_v2.VERSION:
		v.Database = db_v2.Database{}
	case db_v3.VERSION:
		v.Database = db_v3.Database{}
	default:
		return false, errors.Format("unsupported vault version \"%d\"", v.Meta.DatabaseVersion)
	}

	return true, nil
}

func detectMetadata(path string) meta.Encoder {
	_, err := os.Stat(meta_encoder.Encoder{}.MetadataPath(path))
	if err == nil {
		return meta_encoder.Encoder{}
	}

	return nil
}

func loadFromDatabase(v *Vault) (bool, error) {
	db := detectDatabase(v.Path)
	if db == nil {
		return false, nil
	}

	v.Database = db
	v.Meta.DatabaseVersion = db.GetVersion()

	err := v.SaveMetadata()
	if err != nil {
		return false, err
	}

	return true, nil
}

func detectDatabase(path string) database.Database {
	_, err := os.Stat(db_v2.Database{}.IndexPath(path))
	if err == nil {
		return db_v2.Database{}
	}

	_, err = os.Stat(db_v1.Database{}.IndexPath(path))
	if err == nil {
		return db_v1.Database{}
	}

	return nil
}
