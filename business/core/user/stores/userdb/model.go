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

func toDBUser(e user.User) dbUser {
	roles := make([]string, len(e.Roles))
	for i, role := range e.Roles {
		roles[i] = role.Name()
	}

	return dbUser{
		ID:           e.ID,
		FirstName:    e.FirstName,
		LastName:     e.LastName,
		Phone:        e.Phone,
		PasswordHash: e.PasswordHash,
		Roles:        roles,
		Status:       e.Status.Name(),
		CreatedAt:    e.CreatedAt.UTC(),
		UpdatedAt:    e.UpdatedAt.UTC(),
	}
}

func toCoreUser(dbUsr dbUser) (user.User, error) {
	roles := make([]user.Role, len(dbUsr.Roles))
	for i, value := range dbUsr.Roles {
		var err error
		roles[i], err = user.ParseRole(value)
		if err != nil {
			return user.User{}, fmt.Errorf("parse role: %w", err)
		}
	}
	status, err := user.ParseStatus(dbUsr.Status)
	if err != nil {
		return user.User{}, fmt.Errorf("parse status: %w", err)
	}

	usr := user.User{
		ID:           dbUsr.ID,
		FirstName:    dbUsr.FirstName,
		LastName:     dbUsr.LastName,
		Phone:        dbUsr.Phone,
		Roles:        roles,
		PasswordHash: dbUsr.PasswordHash,
		Status:       status,
		CreatedAt:    dbUsr.CreatedAt.In(time.Local),
		UpdatedAt:    dbUsr.UpdatedAt.In(time.Local),
	}

	return usr, nil
}

func toCoreUsers(rows []dbUser) ([]user.User, error) {
	usrs := make([]user.User, len(rows))
	for i, dbUsr := range rows {
		var err error
		usrs[i], err = toCoreUser(dbUsr)
		if err != nil {
			return nil, err
		}
	}
	return usrs, nil
}
