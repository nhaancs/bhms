package user

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

// EventSource represents the source of the given event.
const EventSource = "user"

// Set of user relatated events.
const (
	EventUpdated = "UserUpdated"
)

// =============================================================================

// EventParamsUpdated is the event parameters for the updated event.
type EventParamsUpdated struct {
	UserID uuid.UUID
	UpdateUser
}

// String returns a string representation of the event parameters.
func (p *EventParamsUpdated) String() string {
	return fmt.Sprintf("&EventParamsUpdated{UserID:%v, Enabled:%v}", p.UserID, p.Enabled)
}

// Marshal returns the event parameters encoded as JSON.
func (p *EventParamsUpdated) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

// UnmarshalUpdated parses the event parameters from JSON.
func UnmarshalUpdated(rawParams []byte) (*EventParamsUpdated, error) {
	var params EventParamsUpdated
	err := json.Unmarshal(rawParams, &params)
	if err != nil {
		return nil, fmt.Errorf("expected an encoded %T: %w", params, err)
	}

	return &params, nil
}
