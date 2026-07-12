package record

import (
	"github.com/binarysoupdev/go-commando/errors"

	"github.com/google/uuid"
)

// Version 2
type Record struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`

	Username string         `json:"username"`
	Password string         `json:"password"`
	Other    map[string]any `json:"other"`
}

func NewFromName(name string) Record {
	return Record{
		ID:       uuid.New(),
		Name:     name,
		Username: "",
		Password: "",
		Other:    map[string]interface{}{},
	}
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
