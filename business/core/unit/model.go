package unit

import (
	"github.com/google/uuid"
	"time"
)

type Unit struct {
	ID         uuid.UUID
	Name       string
	PropertyID uuid.UUID
	BlockID    uuid.UUID
	FloorID    uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type NewUnit struct {
	Name       string
	PropertyID uuid.UUID
	BlockID    uuid.UUID
	FloorID    uuid.UUID
}

type UpdateUnit struct {
	Name *string
}
