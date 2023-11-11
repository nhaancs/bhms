package property

import "fmt"

// Set of possible status for a user.
var (
	StatusCreated = Status{"CREATED"}
	StatusActive  = Status{"ACTIVE"}
	StatusDeleted = Status{"DELETED"}
)

// Set of known roles.
var status = map[string]Status{
	StatusCreated.name: StatusCreated,
	StatusActive.name:  StatusActive,
	StatusDeleted.name: StatusDeleted,
}

// Status represents a user status in the system.
type Status struct {
	name string
}

// ParseStatus parses the string value and returns a status if one exists.
func ParseStatus(value string) (Status, error) {
	sts, exists := status[value]
	if !exists {
		return Status{}, fmt.Errorf("invalid role %q", value)
	}

	return sts, nil
}

// MustParseStatus parses the string value and returns a status if one exists. If
// an error occurs the function panics.
func MustParseStatus(value string) Status {
	sts, err := ParseStatus(value)
	if err != nil {
		panic(err)
	}

	return sts
}

// Name returns the name of the status.
func (r *Status) Name() string {
	return r.name
}

// UnmarshalText implement the unmarshal interface for JSON conversions.
func (r *Status) UnmarshalText(data []byte) error {
	role, err := ParseStatus(string(data))
	if err != nil {
		return err
	}

	r.name = role.name
	return nil
}

// MarshalText implement the marshal interface for JSON conversions.
func (r *Status) MarshalText() ([]byte, error) {
	return []byte(r.name), nil
}

// Equal provides support for the go-cmp package and testing.
func (r *Status) Equal(r2 Status) bool {
	return r.name == r2.name
}
