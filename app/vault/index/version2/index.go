package v2

import (
	"path/filepath"

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

func (Index) GetVersion() int {
	return VERSION
}

func (idx Index) Filepath() string {
	return filepath.Join(idx.Path, FILENAME)
}

func (idx Index) RecordPath(id uuid.UUID) string {
	return filepath.Join(idx.Path, id.String())
}

func (idx Index) Upgrade() (Index, error) {
	return idx, nil
}
