package id

import (
	"github.com/google/uuid"
)

type UUID = uuid.UUID

func NewUUID() UUID {
	return UUID(uuid.New())
}

func Parse(s string) (UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return UUID{}, err
	}
	return UUID(id), nil
}
