package block

import (
	"github.com/google/uuid"
	"time"
)

type Block struct {
	ID         uuid.UUID
	Name       string
	PropertyID uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type NewBlock struct {
	Name       string
	PropertyID uuid.UUID
}

type UpdateBlock struct {
	Name *string
}
