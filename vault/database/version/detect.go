package version

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"pvault/vault/database"
	v1 "pvault/vault/database/version/v1"
	v2 "pvault/vault/database/version/v2"

	"github.com/binarysoupdev/go-commando/errors"
)

func Detect(path string) (database.Database, error) {
	_, err := os.Stat(filepath.Join(path, v2.INDEX_FILE))
	if err == nil {
		return detectFromVersionHeader(path, v2.INDEX_FILE)
	}

	_, err = os.Stat(filepath.Join(path, v1.INDEX_FILE))
	if err == nil {
		return v1.New(path), nil
	}

	return nil, errors.New("index file not found")
}

func detectFromVersionHeader(path, indexFile string) (database.Database, error) {
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
		return v2.New(path), nil
	default:
		return nil, errors.Format("unsupported version \"%d\"", version)
	}
}
