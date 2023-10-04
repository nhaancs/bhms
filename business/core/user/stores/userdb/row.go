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

func toUserRow(usr user.UserEntity) userRow {
	roles := make([]string, len(usr.Roles))
	for i, role := range usr.Roles {
		roles[i] = role.Name()
	}

	return userRow{
		ID:           usr.ID,
		Name:         usr.Name,
		Email:        usr.Email.Address,
		Roles:        roles,
		PasswordHash: usr.PasswordHash,
		Department: sql.NullString{
			String: usr.Department,
			Valid:  usr.Department != "",
		},
		Enabled:     usr.Enabled,
		DateCreated: usr.DateCreated.UTC(),
		DateUpdated: usr.DateUpdated.UTC(),
	}
}

func toUserEntity(dbUsr userRow) (user.UserEntity, error) {
	addr := mail.Address{
		Address: dbUsr.Email,
	}

	roles := make([]user.Role, len(dbUsr.Roles))
	for i, value := range dbUsr.Roles {
		var err error
		roles[i], err = user.ParseRole(value)
		if err != nil {
			return user.UserEntity{}, fmt.Errorf("parse role: %w", err)
		}
	}

	usr := user.UserEntity{
		ID:           dbUsr.ID,
		Name:         dbUsr.Name,
		Email:        addr,
		Roles:        roles,
		PasswordHash: dbUsr.PasswordHash,
		Enabled:      dbUsr.Enabled,
		Department:   dbUsr.Department.String,
		DateCreated:  dbUsr.DateCreated.In(time.Local),
		DateUpdated:  dbUsr.DateUpdated.In(time.Local),
	}

	return usr, nil
}

func toUserEntities(dbUsers []userRow) ([]user.UserEntity, error) {
	usrs := make([]user.UserEntity, len(dbUsers))
	for i, dbUsr := range dbUsers {
		var err error
		usrs[i], err = toUserEntity(dbUsr)
		if err != nil {
			return nil, err
		}
	}
	return usrs, nil
}
