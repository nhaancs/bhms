package floor

import (
	"github.com/google/uuid"
	"time"
)

type Floor struct {
	ID         uuid.UUID
	Name       string
	PropertyID uuid.UUID
	BlockID    uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type NewFloor struct {
	Name       string
	PropertyID uuid.UUID
	BlockID    uuid.UUID
}

type UpdateFloor struct {
	Name *string
}
