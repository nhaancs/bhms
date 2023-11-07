package userdb

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nhaancs/bhms/business/core/user"
	"github.com/nhaancs/bhms/business/data/dbsql/pgx/dbarray"
)

// dbUser represent the structure we need for moving data
// between the app and the database.
type dbUser struct {
	ID           uuid.UUID      `db:"id"`
	FirstName    string         `db:"first_name"`
	LastName     string         `db:"last_name"`
	Phone        string         `db:"phone"`
	PasswordHash []byte         `db:"password_hash"`
	Roles        dbarray.String `db:"roles"`
	Status       string         `db:"status"`
	CreatedAt    time.Time      `db:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at"`
}

func toDBUser(c user.User) dbUser {
	roles := make([]string, len(c.Roles))
	for i, role := range c.Roles {
		roles[i] = role.Name()
	}

	return dbUser{
		ID:           c.ID,
		FirstName:    c.FirstName,
		LastName:     c.LastName,
		Phone:        c.Phone,
		PasswordHash: c.PasswordHash,
		Roles:        roles,
		Status:       c.Status.Name(),
		CreatedAt:    c.CreatedAt.UTC(),
		UpdatedAt:    c.UpdatedAt.UTC(),
	}
}

func toCoreUser(r dbUser) (user.User, error) {
	roles := make([]user.Role, len(r.Roles))
	for i, value := range r.Roles {
		var err error
		roles[i], err = user.ParseRole(value)
		if err != nil {
			return user.User{}, fmt.Errorf("parse role: %w", err)
		}
	}
	status, err := user.ParseStatus(r.Status)
	if err != nil {
		return user.User{}, fmt.Errorf("parse status: %w", err)
	}

	usr := user.User{
		ID:           r.ID,
		FirstName:    r.FirstName,
		LastName:     r.LastName,
		Phone:        r.Phone,
		Roles:        roles,
		PasswordHash: r.PasswordHash,
		Status:       status,
		CreatedAt:    r.CreatedAt.In(time.Local),
		UpdatedAt:    r.UpdatedAt.In(time.Local),
	}

	return usr, nil
}

func toCoreUsers(rs []dbUser) ([]user.User, error) {
	usrs := make([]user.User, len(rs))
	var err error
	for i, r := range rs {
		usrs[i], err = toCoreUser(r)
		if err != nil {
			return nil, err
		}
	}
	return usrs, nil
}
