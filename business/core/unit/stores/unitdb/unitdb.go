package unitdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/nhaancs/bhms/business/core/unit"
	db "github.com/nhaancs/bhms/business/data/dbsql/pgx"
	"github.com/nhaancs/bhms/business/data/transaction"
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

func (s *Store) Create(ctx context.Context, core unit.Unit) error {
	const q = `
	INSERT INTO units
		(id, name, property_id, block_id, floor_id, status, created_at, updated_at)
	VALUES
		(:id, :name, :property_id, :block_id, :floor_id, :status, :created_at, :updated_at)`

	if err := db.NamedExecContext(ctx, s.log, s.db, q, toDBUnit(core)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) BatchCreate(ctx context.Context, cores []unit.Unit) error {
	const q = `
	INSERT INTO units
		(id, name, property_id, block_id, floor_id, status, created_at, updated_at)
	VALUES
		(:id, :name, :property_id, :block_id, :floor_id, :status, :created_at, :updated_at)`

	if err := db.NamedExecContext(ctx, s.log, s.db, q, toDBUnits(cores)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) Update(ctx context.Context, core unit.Unit) error {
	const q = `
	UPDATE
		units
	SET 
		"name" = :name,
		"status" = :status,
		"updated_at" = :updated_at
	WHERE
		id = :id`

	if err := db.NamedExecContext(ctx, s.log, s.db, q, toDBUnit(core)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) QueryByID(ctx context.Context, id uuid.UUID) (unit.Unit, error) {
	data := struct {
		ID string `db:"id"`
	}{
		ID: id.String(),
	}

	const q = `
	SELECT
        id, name, property_id, block_id, floor_id, status, created_at, updated_at
	FROM
		units
	WHERE 
		id = :id`

	var row dbUnit
	if err := db.NamedQueryStruct(ctx, s.log, s.db, q, data, &row); err != nil {
		if errors.Is(err, db.ErrDBNotFound) {
			return unit.Unit{}, fmt.Errorf("namedquerystruct: %w", unit.ErrNotFound)
		}
		return unit.Unit{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	unt, err := toCoreUnit(row)
	if err != nil {
		return unit.Unit{}, err
	}

	return unt, nil
}

func (s *Store) QueryByFloorID(ctx context.Context, id uuid.UUID) ([]unit.Unit, error) {
	data := struct {
		FloorID string `db:"floor_id"`
	}{
		FloorID: id.String(),
	}

	const q = `
	SELECT
        id, name, property_id, block_id, floor_id, status, created_at, updated_at
	FROM
		units
	WHERE
		floor_id = :floor_id`

	var rows []dbUnit
	if err := db.NamedQuerySlice(ctx, s.log, s.db, q, data, &rows); err != nil {
		if errors.Is(err, db.ErrDBNotFound) {
			return nil, fmt.Errorf("namedqueryslice: %w", unit.ErrNotFound)
		}
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	unts, err := toCoreUnits(rows)
	if err != nil {
		return nil, err
	}

	return unts, nil
}

func (s *Store) QueryByPropertyID(ctx context.Context, id uuid.UUID) ([]unit.Unit, error) {
	data := struct {
		PropertyID string `db:"property_id"`
	}{
		PropertyID: id.String(),
	}

	const q = `
	SELECT
        id, name, property_id, block_id, floor_id, status, created_at, updated_at
	FROM
		units
	WHERE
		property_id = :property_id`

	var rows []dbUnit
	if err := db.NamedQuerySlice(ctx, s.log, s.db, q, data, &rows); err != nil {
		if errors.Is(err, db.ErrDBNotFound) {
			return nil, fmt.Errorf("namedqueryslice: %w", unit.ErrNotFound)
		}
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	unts, err := toCoreUnits(rows)
	if err != nil {
		return nil, err
	}

	return unts, nil
}

// ExecuteUnderTransaction constructs a new Store value replacing the sqlx DB
// value with a sqlx DB value that is currently inside a transaction.
func (s *Store) ExecuteUnderTransaction(tx transaction.Transaction) (unit.Storer, error) {
	ec, err := db.GetExtContext(tx)
	if err != nil {
		return nil, err
	}

	s = &Store{
		log: s.log,
		db:  ec,
	}

	return s, nil
}
