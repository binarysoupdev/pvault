package meta

const VERSION = 1

type Metadata struct {
	DatabaseVersion int
}

func New(dbVersion int) Metadata {
	return Metadata{
		DatabaseVersion: dbVersion,
	}
}
