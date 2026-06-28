package vault

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"pvault/errors"
	"pvault/vault/data"
	"pvault/vault/data/legacy"
)

const (
	LEGACY_INDEX_FILE = "index.txt"
	LEGACY_VERSION    = 1
)

func Open(path string) (Vault, error) {
	v := Vault{
		Path: path,
	}

	version, err := v.detectVersion()
	if err != nil {
		return Vault{}, err
	}

	v.Version = version
	switch version {
	case 1:
		v.Database = legacy.DatabaseV1{}
	case 2:
		v.Database = data.DatabaseV2{}
	default:
		return Vault{}, errors.Format("unsupported version \"%d\"", version)
	}

	v.Index, err = v.Database.LoadIndex(v.IndexPath())
	if err != nil {
		return v, errors.Chain(err, "error parsing index file")
	}

	return v, nil
}

func (v Vault) detectVersion() (uint16, error) {
	_, err := os.Stat(v.IndexPath())
	if err == nil {
		return v.parseVersionHeader()
	}

	_, err = os.Stat(filepath.Join(v.Path, LEGACY_INDEX_FILE))
	if err == nil {
		return LEGACY_VERSION, nil
	}

	return 0, errors.New("index file not found")
}

func (v Vault) parseVersionHeader() (uint16, error) {
	header := make([]byte, 2)

	file, err := os.Open(v.IndexPath())
	if err != nil {
		return 0, errors.Chain(err, "error opening index file")
	}
	defer file.Close()

	_, err = file.Read(header)
	if err != nil {
		return 0, errors.Chain(err, "error reading version header")
	}

	return binary.BigEndian.Uint16(header), nil
}
