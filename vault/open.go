package vault

import (
	"os"
	"path/filepath"
	"pvault/vault/database"
	db_v1 "pvault/vault/database/encoder/legacy/v1"
	db_v2 "pvault/vault/database/encoder/legacy/v2"
	db_v3 "pvault/vault/database/encoder/v3"
	"pvault/vault/meta"
	meta_v1 "pvault/vault/meta/encoder/v1"

	"github.com/binarysoupdev/go-commando/errors"
)

type ErrorNotFound struct{}

func (err ErrorNotFound) Error() string {
	return "vault not found"
}

//=================================

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

	return ErrorNotFound{}
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
		v.DatabaseEncoder = db_v1.Encoder{}
	case db_v2.VERSION:
		v.DatabaseEncoder = db_v2.Encoder{}
	case db_v3.VERSION:
		v.DatabaseEncoder = db_v3.Encoder{}
	default:
		return false, errors.Format("unsupported vault version \"%d\"", v.Meta.DatabaseVersion)
	}

	return true, nil
}

func detectMetadata(path string) meta.Encoder {
	_, err := os.Stat(meta_v1.Encoder{}.MetadataPath(path))
	if err == nil {
		return meta_v1.Encoder{}
	}

	return nil
}

func loadFromDatabase(v *Vault) (bool, error) {
	db := detectDatabase(v.Path)
	if db == nil {
		return false, nil
	}

	v.DatabaseEncoder = db
	v.Meta.DatabaseVersion = db.GetVersion()

	err := v.SaveMetadata()
	if err != nil {
		return false, err
	}

	return true, nil
}

func detectDatabase(path string) database.Encoder {
	_, err := os.Stat(db_v2.Encoder{}.IndexPath(path))
	if err == nil {
		return db_v2.Encoder{}
	}

	_, err = os.Stat(db_v1.Encoder{}.IndexPath(path))
	if err == nil {
		return db_v1.Encoder{}
	}

	return nil
}
