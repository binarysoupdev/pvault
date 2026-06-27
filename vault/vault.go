package vault

import (
	"os"
	"pvault/errors"
	"pvault/vault/index"
)

type Vault struct {
	Version int
	Path    string
	Index   index.IndexMap
}

func InitializeNew(path string) (Vault, error) {
	v := Vault{
		Path:  path,
		Index: index.IndexMap{},
	}

	_, err := os.Stat(path)
	if err == nil || !os.IsNotExist(err) {
		return v, errors.New("vault path already exists")
	}

	err = os.MkdirAll(path, 0755)
	if err != nil {
		return v, errors.Chain(err, "error creating vault directory")
	}

	err = index.Codec{}.Encode(v.IndexPath(), v.Index)
	if err != nil {
		return v, errors.Chain(err, "error saving index file")
	}

	return v, nil
}

func Open(path string) (Vault, error) {
	v := Vault{
		Version: index.CURRENT_VERSION,
		Path:    path,
	}

	decoder, err := v.selectDecoder(v.Path)
	if err != nil {
		return Vault{}, err
	}

	if decoder.Version() < v.Version {
		return Vault{}, newOutOfDateError(decoder.Version())
	}

	v.Index, err = decoder.Decode(v.IndexPath())
	if err != nil {
		return v, errors.Chain(err, "error parsing index file")
	}

	return v, nil
}
