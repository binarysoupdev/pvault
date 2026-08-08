package meta

import "time"

const VERSION = 1

type Metadata struct {
	DatabaseVersion int
	Nickname        string
	CreationDate    time.Time
}

func New(dbVersion int, name string) Metadata {
	return Metadata{
		DatabaseVersion: dbVersion,
		Nickname:        name,
		CreationDate:    time.Now(),
	}
}
