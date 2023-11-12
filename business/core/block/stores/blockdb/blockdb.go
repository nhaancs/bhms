package blockdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/nhaancs/bhms/business/core/block"
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

func (s *Store) Create(ctx context.Context, core block.Block) error {
	const q = `
	INSERT INTO blocks
		(id, name, property_id, status, created_at, updated_at)
	VALUES
		(:id, :name, :property_id, :status, :created_at, :updated_at)`

	if err := db.NamedExecContext(ctx, s.log, s.db, q, toDBBlock(core)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) BatchCreate(ctx context.Context, cores []block.Block) error {
	const q = `
	INSERT INTO blocks
		(id, name, property_id, status, created_at, updated_at)
	VALUES
		(:id, :name, :property_id, :status, :created_at, :updated_at)`

	if err := db.NamedExecContext(ctx, s.log, s.db, q, toDBBlocks(cores)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) Update(ctx context.Context, core block.Block) error {
	const q = `
	UPDATE
		blocks
	SET 
		"name" = :name,
		"status" = :status,
		"updated_at" = :updated_at
	WHERE
		id = :id AND status != 'DELETED'`

	if err := db.NamedExecContext(ctx, s.log, s.db, q, toDBBlock(core)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) QueryByID(ctx context.Context, id uuid.UUID) (block.Block, error) {
	data := struct {
		ID     string `db:"id"`
		Status string `db:"status"`
	}{
		ID:     id.String(),
		Status: block.StatusDeleted.Name(),
	}

	const q = `
	SELECT
        id, name, property_id, status, created_at, updated_at
	FROM
		blocks
	WHERE 
		id = :id AND status != :status`

	var row dbBlock
	if err := db.NamedQueryStruct(ctx, s.log, s.db, q, data, &row); err != nil {
		if errors.Is(err, db.ErrDBNotFound) {
			return block.Block{}, fmt.Errorf("namedquerystruct: %w", block.ErrNotFound)
		}
		return block.Block{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	blck, err := toCoreBlock(row)
	if err != nil {
		return block.Block{}, err
	}

	return blck, nil
}

func (s *Store) QueryByPropertyID(ctx context.Context, id uuid.UUID) ([]block.Block, error) {
	data := struct {
		PropertyID string `db:"property_id"`
		Status     string `db:"status"`
	}{
		PropertyID: id.String(),
		Status:     block.StatusDeleted.Name(),
	}

	const q = `
	SELECT
        id, name, property_id, status, created_at, updated_at
	FROM
		blocks
	WHERE
		property_id = :property_id AND status != :status`

	var rows []dbBlock
	if err := db.NamedQuerySlice(ctx, s.log, s.db, q, data, &rows); err != nil {
		if errors.Is(err, db.ErrDBNotFound) {
			return nil, fmt.Errorf("namedqueryslice: %w", block.ErrNotFound)
		}
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	blcks, err := toCoreBlocks(rows)
	if err != nil {
		return nil, err
	}

	return blcks, nil
}

// ExecuteUnderTransaction constructs a new Store value replacing the sqlx DB
// value with a sqlx DB value that is currently inside a transaction.
func (s *Store) ExecuteUnderTransaction(tx transaction.Transaction) (block.Storer, error) {

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
