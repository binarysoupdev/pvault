package database

import (
	"encoding/binary"
	"os"
	"path/filepath"
	v1 "pvault/vault/database/version1"
	v2 "pvault/vault/database/version2"
	"pvault/vault/index"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
)

type Index interface {
	GetVersion() int

	Filepath() string
	RecordPath(id uuid.UUID) string

	SaveIndex(idx index.IndexMap) error
	LoadIndex() (index.IndexMap, error)
	Upgrade(idx index.IndexMap) error
}

func detectIndex(path string) (Index, error) {
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
	header := make([]byte, 2)

	file, err := os.Open(filepath.Join(path, filename))
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
		return v2.NewIndex(path), nil
	default:
		return nil, errors.Format("unsupported version \"%d\"", version)
	}
}
