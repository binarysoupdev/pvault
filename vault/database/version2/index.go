package v2

import (
	"path/filepath"
	"pvault/vault/index"

	"github.com/google/uuid"
)

const (
	VERSION  = 2
	FILENAME = "index.bin"
)

type Index struct {
	Path string
}

func NewIndex(path string) Index {
	return Index{
		Path: path,
	}
}

func (idx Index) Filepath() string {
	return filepath.Join(idx.Path, FILENAME)
}

func (Index) GetVersion() int {
	return VERSION
}

func (idx Index) RecordPath(id uuid.UUID) string {
	return filepath.Join(idx.Path, id.String())
}

func (Index) Upgrade(m index.IndexMap) error {
	return nil
}
