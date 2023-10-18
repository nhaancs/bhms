package sms

import "fmt"

// Type represents an sms type.
type Type struct {
	id string
}

// ID returns the ID of the role.
func (t *Type) ID() string {
	return t.id
}

// Set of possible sms types.
var (
	// TypeBrandName Brand name chăm sóc khách hàng
	TypeBrandName = Type{"2"}
	// Type10Numbers Tin nhắn đầu số cố định 10 số, chuyên dùng cho chăm sóc khách hàng.
	Type10Numbers = Type{"8"}
	// TypePriorityZalo Tin nhắn Zalo ưu tiên
	TypePriorityZalo = Type{"24"}
	// TypeNormalZalo Tin nhắn Zalo thường
	TypeNormalZalo = Type{"25"}
)

// Set of known types.
var types = map[string]Type{
	TypeBrandName.id:    TypeBrandName,
	Type10Numbers.id:    Type10Numbers,
	TypePriorityZalo.id: TypePriorityZalo,
	TypeNormalZalo.id:   TypeNormalZalo,
}

// ParseType parses the string value and returns a type if one exists.
func ParseType(value string) (Type, error) {
	t, exists := types[value]
	if !exists {
		return Type{}, fmt.Errorf("invalid type %q", value)
	}

	return t, nil
}

// UnmarshalText implement the unmarshal interface for JSON conversions.
func (t *Type) UnmarshalText(data []byte) error {
	role, err := ParseType(string(data))
	if err != nil {
		return err
	}

	t.id = role.id
	return nil
}

// MarshalText implement the marshal interface for JSON conversions.
func (t *Type) MarshalText() ([]byte, error) {
	return []byte(t.id), nil
}

// Equal provides support for the go-cmp package and testing.
func (t *Type) Equal(t2 Type) bool {
	return t.id == t2.id
}
