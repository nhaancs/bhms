package userdb

import (
	"database/sql"
	"fmt"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/nhaancs/bhms/business/core/user"
	"github.com/nhaancs/bhms/business/data/dbsql/pgx/dbarray"
)

// userRow represent the structure we need for moving data
// between the app and the database.
type userRow struct {
	ID           uuid.UUID      `db:"user_id"`
	Name         string         `db:"name"`
	Email        string         `db:"email"`
	Roles        dbarray.String `db:"roles"`
	PasswordHash []byte         `db:"password_hash"`
	Enabled      bool           `db:"enabled"`
	Department   sql.NullString `db:"department"`
	DateCreated  time.Time      `db:"date_created"`
	DateUpdated  time.Time      `db:"date_updated"`
}

func toUserRow(e user.UserEntity) userRow {
	roles := make([]string, len(e.Roles))
	for i, role := range e.Roles {
		roles[i] = role.Name()
	}

	return userRow{
		ID:           e.ID,
		Name:         e.Name,
		Email:        e.Email.Address,
		Roles:        roles,
		PasswordHash: e.PasswordHash,
		Department: sql.NullString{
			String: e.Department,
			Valid:  e.Department != "",
		},
		Enabled:     e.Enabled,
		DateCreated: e.DateCreated.UTC(),
		DateUpdated: e.DateUpdated.UTC(),
	}
}

func toUserEntity(r userRow) (user.UserEntity, error) {
	addr := mail.Address{
		Address: r.Email,
	}

	roles := make([]user.Role, len(r.Roles))
	for i, value := range r.Roles {
		var err error
		roles[i], err = user.ParseRole(value)
		if err != nil {
			return user.UserEntity{}, fmt.Errorf("parse role: %w", err)
		}
	}

	usr := user.UserEntity{
		ID:           r.ID,
		Name:         r.Name,
		Email:        addr,
		Roles:        roles,
		PasswordHash: r.PasswordHash,
		Enabled:      r.Enabled,
		Department:   r.Department.String,
		DateCreated:  r.DateCreated.In(time.Local),
		DateUpdated:  r.DateUpdated.In(time.Local),
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
