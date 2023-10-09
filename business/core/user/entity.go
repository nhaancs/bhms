package user

import (
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/nhaancs/bhms/business/core/event"
)

// UserEntity represents information about an individual user.
type UserEntity struct {
	ID           uuid.UUID
	Name         string
	Email        mail.Address
	Roles        []Role
	PasswordHash []byte
	Department   string
	Enabled      bool
	DateCreated  time.Time
	DateUpdated  time.Time
}

// RegisterEntity contains information needed to create a new user.
type RegisterEntity struct {
	Name            string
	Email           mail.Address
	Roles           []Role
	Department      string
	Password        string
	PasswordConfirm string
}

// UpdateUserEntity contains information needed to update a user.
type UpdateUserEntity struct {
	Name            *string
	Email           *mail.Address
	Roles           []Role
	Department      *string
	Password        *string
	PasswordConfirm *string
	Enabled         *bool
}

// UpdatedEvent constructs an event for when a user is updated.
func (uu UpdateUserEntity) UpdatedEvent(userID uuid.UUID) event.Event {
	params := EventParamsUpdated{
		UserID: userID,
		UpdateUserEntity: UpdateUserEntity{
			Enabled: uu.Enabled,
		},
	}

	rawParams, err := params.Marshal()
	if err != nil {
		panic(err)
	}

	return event.Event{
		Source:    EventSource,
		Type:      EventUpdated,
		RawParams: rawParams,
	}
}
