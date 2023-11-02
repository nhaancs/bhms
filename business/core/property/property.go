package property

import (
	"context"
	"github.com/google/uuid"
)

// TODO: limit the number of properties can be created by a user

type Storer interface {
	Create(ctx context.Context, prprty Property) error
	Update(ctx context.Context, prprty Property) error
	Delete(ctx context.Context, prprty Property) error
	QueryByID(ctx context.Context, prprtyID uuid.UUID) (Property, error)
	QueryByManagerID(ctx context.Context, managerID uuid.UUID) ([]Property, error)
}
