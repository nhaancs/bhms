package floordb

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/nhaancs/bhms/business/core/floor"
	db "github.com/nhaancs/bhms/business/data/dbsql/pgx"
	"github.com/nhaancs/bhms/foundation/logger"
)

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

func (s *Store) Create(ctx context.Context, core floor.Floor) error {
	const q = `
	INSERT INTO floors
		(id, name, property_id, block_id, created_at, updated_at)
	VALUES
		(:id, :name, :property_id, :block_id, :created_at, :updated_at)`

	if err := db.NamedExecContext(ctx, s.log, s.db, q, toDBFloor(core)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) BatchCreate(ctx context.Context, cores []floor.Floor) error {
	const q = `
	INSERT INTO floors
		(id, name, property_id, block_id, created_at, updated_at)
	VALUES
		(:id, :name, :property_id, :block_id, :created_at, :updated_at)`

	if err := db.NamedExecContext(ctx, s.log, s.db, q, toDBFloors(cores)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) Update(ctx context.Context, core floor.Floor) error {
	const q = `
	UPDATE
		floors
	SET 
		"name" = :name,
		"updated_at" = :updated_at
	WHERE
		id = :id`

	if err := db.NamedExecContext(ctx, s.log, s.db, q, toDBFloor(core)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) Delete(ctx context.Context, core floor.Floor) error {
	data := struct {
		ID string `db:"id"`
	}{
		ID: core.ID.String(),
	}

	const q = `
	DELETE FROM
		floors
	WHERE
		id = :id`

	if err := db.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) QueryByID(ctx context.Context, id uuid.UUID) (floor.Floor, error) {
	data := struct {
		ID string `db:"id"`
	}{
		ID: id.String(),
	}

	const q = `
	SELECT
        id, name, property_id, block_id, created_at, updated_at
	FROM
		floors
	WHERE 
		id = :id`

	var row dbFloor
	if err := db.NamedQueryStruct(ctx, s.log, s.db, q, data, &row); err != nil {
		if errors.Is(err, db.ErrDBNotFound) {
			return floor.Floor{}, fmt.Errorf("namedquerystruct: %w", floor.ErrNotFound)
		}
		return floor.Floor{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	flr, err := toCoreFloor(row)
	if err != nil {
		return floor.Floor{}, err
	}

	return flr, nil
}

func (s *Store) QueryByBlockID(ctx context.Context, blockID uuid.UUID) ([]floor.Floor, error) {
	data := struct {
		BlockID string `db:"block_id"`
	}{
		BlockID: blockID.String(),
	}

	const q = `
	SELECT
        id, name, property_id, block_id, created_at, updated_at
	FROM
		floors
	WHERE
		block_id = :block_id`

	var rows []dbFloor
	if err := db.NamedQuerySlice(ctx, s.log, s.db, q, data, &rows); err != nil {
		if errors.Is(err, db.ErrDBNotFound) {
			return nil, fmt.Errorf("namedqueryslice: %w", floor.ErrNotFound)
		}
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	flrs, err := toCoreFloors(rows)
	if err != nil {
		return nil, err
	}

	return flrs, nil
}
