package v1

import (
	"path/filepath"

	"github.com/google/uuid"
)

const (
	VERSION  = 1
	FILENAME = "index.txt"
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
	return filepath.Join(idx.Path, id.String()+".crypt")
}
