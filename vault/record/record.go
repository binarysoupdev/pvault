package record

import (
	v2 "pvault/vault/record/record/v2"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
)

type Record interface {
	GetVersion() int
	GetID() uuid.UUID
	GetName() string
	Upgrade() v2.Record
}

func Validate(r Record) error {
	errs := errors.Errors{}

	if r.GetID() == uuid.Nil {
		errs.AddNew("id cannot be nil (all zeroes)")
	}

	if r.GetName() == "" {
		errs.AddNew("name cannot be empty")
	}

	return errs.Collapse(", ")
}
