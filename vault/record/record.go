package record

import (
	"pvault/errors"

	"github.com/google/uuid"
)

type Record struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`

	Username string `json:"username"`
	Password string `json:"password"`
	Other    any    `json:"other"`
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
		errs.Add("\"ID\" cannot be nil (all zeroes)")
	}
	if r.Name == "" {
		errs.Add("\"Name\" cannot be empty")
	}

	return errs.Collapse(", ")
}
