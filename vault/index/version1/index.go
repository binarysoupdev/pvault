package v1

import (
	"path/filepath"

	"github.com/google/uuid"
)

const (
	VERSION    = 1
	INDEX_FILE = "index.txt"
)

type Index struct {
	Path string
}

func NewIndex(path string) Index {
	return Index{
		Path: path,
	}
}

func (idx Index) filepath() string {
	return filepath.Join(idx.Path, INDEX_FILE)
}

func (Index) GetVersion() uint16 {
	return VERSION
}

func (idx Index) RecordPath(id uuid.UUID) string {
	return filepath.Join(idx.Path, id.String()+".crypt")
}
