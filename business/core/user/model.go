package user

import (
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/nhaancs/realworld/business/core/event"
)

// User represents information about an individual user.
type User struct {
	ID           uuid.UUID
	Username     string
	Email        mail.Address
	Bio          string
	Image        string
	PasswordHash []byte
	Enabled      bool
	DateCreated  time.Time
	DateUpdated  time.Time
}

// NewUser contains information needed to create a new user.
type NewUser struct {
	Username string
	Email    mail.Address
	Password string
}

// UpdateUser contains information needed to update a user.
type UpdateUser struct {
	Email   *mail.Address
	Bio     *string
	Image   *string
	Enabled *bool
}

// UpdatedEvent constructs an event for when a user is updated.
func (uu UpdateUser) UpdatedEvent(userID uuid.UUID) event.Event {
	params := EventParamsUpdated{
		UserID: userID,
		UpdateUser: UpdateUser{
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
