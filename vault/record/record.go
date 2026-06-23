package record

import (
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
