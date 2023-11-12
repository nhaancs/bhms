package property

import (
	"github.com/google/uuid"
	"time"
)

type Property struct {
	ID              uuid.UUID
	ManagerID       uuid.UUID
	Name            string
	AddressLevel1ID uint32
	AddressLevel2ID uint32
	AddressLevel3ID uint32
	Street          string
	Status          Status
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type NewProperty struct {
	ID              uuid.UUID
	ManagerID       uuid.UUID
	Name            string
	AddressLevel1ID uint32
	AddressLevel2ID uint32
	AddressLevel3ID uint32
	Street          string
	Status          Status
}

type UpdateProperty struct {
	Name            *string
	AddressLevel1ID *uint32
	AddressLevel2ID *uint32
	AddressLevel3ID *uint32
	Street          *string
}

type UpdatePropertyStatus struct {
	Status *Status
}
