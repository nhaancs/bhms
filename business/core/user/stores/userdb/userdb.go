// Package userdb contains user related CRUD functionality.
package userdb

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/nhaancs/bhms/business/core/user"
	db "github.com/nhaancs/bhms/business/data/dbsql/pgx"
	"github.com/nhaancs/bhms/business/data/dbsql/pgx/dbarray"
	"github.com/nhaancs/bhms/foundation/logger"
)

// Store manages the set of APIs for user database access.
type Store struct {
	log *logger.Logger
	db  sqlx.ExtContext
}

// NewStore constructs the api for data access.
func NewStore(log *logger.Logger, db *sqlx.DB) *Store {
	return &Store{
		log: log,
		db:  db,
	}
}

// Create inserts a new user into the database.
func (s *Store) Create(ctx context.Context, usr user.UserEntity) error {
	const q = `
	INSERT INTO users
		(id, first_name, last_name, phone, password_hash, roles, status, created_at, updated_at)
	VALUES
		(:id, :first_name, :last_name, :phone, :password_hash, :roles, :status, :created_at, :updated_at)`

	if err := db.NamedExecContext(ctx, s.log, s.db, q, toUserRow(usr)); err != nil {
		if errors.Is(err, db.ErrDBDuplicatedEntry) {
			return fmt.Errorf("namedexeccontext: %w", user.ErrUniquePhone)
		}
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Update replaces a user document in the database.
func (s *Store) Update(ctx context.Context, usr user.UserEntity) error {
	const q = `
	UPDATE
		users
	SET 
		"first_name" = :first_name,
		"last_name" = :last_name,
		"phone" = :phone,
		"roles" = :roles,
		"password_hash" = :password_hash,
		"status" = :status,
		"date_updated" = :date_updated
	WHERE
		id = :id`

	if err := db.NamedExecContext(ctx, s.log, s.db, q, toUserRow(usr)); err != nil {
		if errors.Is(err, db.ErrDBDuplicatedEntry) {
			return user.ErrUniquePhone
		}
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Delete removes a user from the database.
func (s *Store) Delete(ctx context.Context, usr user.UserEntity) error {
	data := struct {
		ID string `db:"id"`
	}{
		ID: usr.ID.String(),
	}

	const q = `
	DELETE FROM
		users
	WHERE
		id = :id`

	if err := db.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// QueryByID gets the specified user from the database.
func (s *Store) QueryByID(ctx context.Context, userID uuid.UUID) (user.UserEntity, error) {
	data := struct {
		ID string `db:"id"`
	}{
		ID: userID.String(),
	}

	const q = `
	SELECT
        id, first_name, last_name, phone, password_hash, roles, status, created_at, updated_at
	FROM
		users
	WHERE 
		id = :id`

	var dbUsr userRow
	if err := db.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbUsr); err != nil {
		if errors.Is(err, db.ErrDBNotFound) {
			return user.UserEntity{}, fmt.Errorf("namedquerystruct: %w", user.ErrNotFound)
		}
		return user.UserEntity{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	usr, err := toUserEntity(dbUsr)
	if err != nil {
		return user.UserEntity{}, err
	}

	return usr, nil
}

// QueryByIDs gets the specified users from the database.
func (s *Store) QueryByIDs(ctx context.Context, userIDs []uuid.UUID) ([]user.UserEntity, error) {
	ids := make([]string, len(userIDs))
	for i, userID := range userIDs {
		ids[i] = userID.String()
	}

	data := struct {
		ID interface {
			driver.Valuer
			sql.Scanner
		} `db:"id"`
	}{
		ID: dbarray.Array(ids),
	}

	const q = `
	SELECT
        id, first_name, last_name, phone, password_hash, roles, status, created_at, updated_at
	FROM
		users
	WHERE
		id = ANY(:id)`

	var rows []userRow
	if err := db.NamedQuerySlice(ctx, s.log, s.db, q, data, &rows); err != nil {
		if errors.Is(err, db.ErrDBNotFound) {
			return nil, user.ErrNotFound
		}
		return nil, fmt.Errorf("namedquerystruct: %w", err)
	}

	usrs, err := toUserEntities(rows)
	if err != nil {
		return nil, err
	}

	return usrs, nil
}

// QueryByPhone gets the specified user from the database by email.
func (s *Store) QueryByPhone(ctx context.Context, phone string) (user.UserEntity, error) {
	data := struct {
		Phone string `db:"phone"`
	}{
		Phone: phone,
	}

	const q = `
	SELECT
        id, first_name, last_name, phone, password_hash, roles, status, created_at, updated_at
	FROM
		users
	WHERE
		phone = :phone`

	var row userRow
	if err := db.NamedQueryStruct(ctx, s.log, s.db, q, data, &row); err != nil {
		if errors.Is(err, db.ErrDBNotFound) {
			return user.UserEntity{}, fmt.Errorf("namedquerystruct: %w", user.ErrNotFound)
		}
		return user.UserEntity{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	usr, err := toUserEntity(row)
	if err != nil {
		return user.UserEntity{}, err
	}

	return usr, nil
}
