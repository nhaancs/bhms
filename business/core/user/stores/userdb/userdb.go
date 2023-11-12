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
func (s *Store) Create(ctx context.Context, core user.User) error {
	const q = `
	INSERT INTO users
		(id, first_name, last_name, phone, password_hash, roles, status, created_at, updated_at)
	VALUES
		(:id, :first_name, :last_name, :phone, :password_hash, :roles, :status, :created_at, :updated_at)`

	if err := db.NamedExecContext(ctx, s.log, s.db, q, toDBUser(core)); err != nil {
		if errors.Is(err, db.ErrDBDuplicatedEntry) {
			return fmt.Errorf("namedexeccontext: %w", user.ErrUniquePhone)
		}
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Update replaces a user document in the database.
func (s *Store) Update(ctx context.Context, core user.User) error {
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
		"updated_at" = :updated_at
	WHERE
		id = :id AND status != 'DELETED'`

	if err := db.NamedExecContext(ctx, s.log, s.db, q, toDBUser(core)); err != nil {
		if errors.Is(err, db.ErrDBDuplicatedEntry) {
			return user.ErrUniquePhone
		}
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// QueryByID gets the specified user from the database.
func (s *Store) QueryByID(ctx context.Context, id uuid.UUID) (user.User, error) {
	data := struct {
		ID     string `db:"id"`
		Status string `db:"status"`
	}{
		ID:     id.String(),
		Status: user.StatusDeleted.Name(),
	}

	const q = `
	SELECT
        id, first_name, last_name, phone, password_hash, roles, status, created_at, updated_at
	FROM
		users
	WHERE 
		id = :id AND status != :status`

	var row dbUser
	if err := db.NamedQueryStruct(ctx, s.log, s.db, q, data, &row); err != nil {
		if errors.Is(err, db.ErrDBNotFound) {
			return user.User{}, fmt.Errorf("namedquerystruct: %w", user.ErrNotFound)
		}
		return user.User{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	usr, err := toCoreUser(row)
	if err != nil {
		return user.User{}, err
	}

	return usr, nil
}

// QueryByIDs gets the specified users from the database.
func (s *Store) QueryByIDs(ctx context.Context, ids []uuid.UUID) ([]user.User, error) {
	uIDs := make([]string, len(ids))
	for i, userID := range ids {
		uIDs[i] = userID.String()
	}

	data := struct {
		ID interface {
			driver.Valuer
			sql.Scanner
		} `db:"id"`
		Status string `db:"status"`
	}{
		ID:     dbarray.Array(uIDs),
		Status: user.StatusDeleted.Name(),
	}

	const q = `
	SELECT
        id, first_name, last_name, phone, password_hash, roles, status, created_at, updated_at
	FROM
		users
	WHERE
		id = ANY(:id) AND status != :status`

	var rows []dbUser
	if err := db.NamedQuerySlice(ctx, s.log, s.db, q, data, &rows); err != nil {
		if errors.Is(err, db.ErrDBNotFound) {
			return nil, user.ErrNotFound
		}
		return nil, fmt.Errorf("namedquerystruct: %w", err)
	}

	usrs, err := toCoreUsers(rows)
	if err != nil {
		return nil, err
	}

	return usrs, nil
}

// QueryByPhone gets the specified user from the database by email.
func (s *Store) QueryByPhone(ctx context.Context, phone string) (user.User, error) {
	data := struct {
		Phone  string `db:"phone"`
		Status string `db:"status"`
	}{
		Phone:  phone,
		Status: user.StatusDeleted.Name(),
	}

	const q = `
	SELECT
        id, first_name, last_name, phone, password_hash, roles, status, created_at, updated_at
	FROM
		users
	WHERE
		phone = :phone AND status != :status`

	var row dbUser
	if err := db.NamedQueryStruct(ctx, s.log, s.db, q, data, &row); err != nil {
		if errors.Is(err, db.ErrDBNotFound) {
			return user.User{}, fmt.Errorf("namedquerystruct: %w", user.ErrNotFound)
		}
		return user.User{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	usr, err := toCoreUser(row)
	if err != nil {
		return user.User{}, err
	}

	return usr, nil
}
