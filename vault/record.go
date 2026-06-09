package vault

import "github.com/google/uuid"

type Record struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`

	Username string `json:"username"`
	Password string `json:"password"`
	Other    any    `json:"other"`
}

func NewRecord(name string) Record {
	return Record{
		ID:   uuid.New(),
		Name: name,
	}
}
