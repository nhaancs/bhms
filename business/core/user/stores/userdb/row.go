package userdb

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nhaancs/bhms/business/core/user"
	"github.com/nhaancs/bhms/business/data/dbsql/pgx/dbarray"
)

// userRow represent the structure we need for moving data
// between the app and the database.
type userRow struct {
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

func toUserRow(e user.UserEntity) userRow {
	roles := make([]string, len(e.Roles))
	for i, role := range e.Roles {
		roles[i] = role.Name()
	}

	return userRow{
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

func toUserEntity(r userRow) (user.UserEntity, error) {
	roles := make([]user.Role, len(r.Roles))
	for i, value := range r.Roles {
		var err error
		roles[i], err = user.ParseRole(value)
		if err != nil {
			return user.UserEntity{}, fmt.Errorf("parse role: %w", err)
		}
	}
	status, err := user.ParseStatus(r.Status)
	if err != nil {
		return user.UserEntity{}, fmt.Errorf("parse status: %w", err)
	}

	usr := user.UserEntity{
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

func toUserEntities(rows []userRow) ([]user.UserEntity, error) {
	usrs := make([]user.UserEntity, len(rows))
	for i, dbUsr := range rows {
		var err error
		usrs[i], err = toUserEntity(dbUsr)
		if err != nil {
			return nil, err
		}
	}
	return usrs, nil
}
