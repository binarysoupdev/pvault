package v2

import (
	"github.com/binarysoupdev/go-commando/errors"

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

func (r Record) GetID() uuid.UUID {
	return r.ID
}

func (r Record) GetName() string {
	return r.Name
}

func (r Record) Validate() error {
	errs := errors.Errors{}

	if r.ID == uuid.Nil {
		errs.AddNew("\"ID\" cannot be nil (all zeroes)")
	}
	if r.Name == "" {
		errs.AddNew("\"Name\" cannot be empty")
	}

	return errs.Collapse(", ")
}

func (r Record) Upgrade() Record {
	return r
}
