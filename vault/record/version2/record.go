package v2

import (
	"github.com/google/uuid"
)

const VERSION = 2

type Record struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`

	Username string         `json:"username"`
	Password string         `json:"password"`
	Other    map[string]any `json:"other"`
}

func NewEmptyRecord(name string) Record {
	return Record{
		ID:       uuid.New(),
		Name:     name,
		Username: "",
		Password: "",
		Other:    map[string]interface{}{},
	}
}

func (r Record) GetVersion() int {
	return VERSION
}

func (r Record) GetID() uuid.UUID {
	return r.ID
}

func (r Record) GetName() string {
	return r.Name
}

func (r Record) Upgrade() Record {
	return r
}
