package index

import (
	"encoding/binary"
	"os"
	"path/filepath"
	v1 "pvault/vault/index/version1"
	v2 "pvault/vault/index/version2"

	"github.com/binarysoupdev/go-commando/errors"
)

func Load(path string) (Index, error) {
	_, err := os.Stat(filepath.Join(path, v2.FILENAME))
	if err == nil {
		return detectFromVersionHeader(path, v2.FILENAME)
	}

	_, err = os.Stat(filepath.Join(path, v1.FILENAME))
	if err == nil {
		return v1.NewIndex(path), nil
	}

	return nil, errors.New("index file not found")
}

func detectFromVersionHeader(path, filename string) (Index, error) {
	bytes, err := os.ReadFile(filepath.Join(path, filename))
	if err != nil {
		return nil, errors.Chain(err, "error reading index file")
	}

	version := binary.BigEndian.Uint16(bytes)

	switch version {
	case 2:
		return v2.NewIndex(path), nil
	default:
		return nil, errors.Format("unsupported index version \"%d\"", version)
	}
}
