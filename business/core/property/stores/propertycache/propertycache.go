package propertycache

import (
	"context"
	"github.com/google/uuid"
	"github.com/nhaancs/bhms/business/core/property"
	"github.com/nhaancs/bhms/foundation/logger"
	"sync"
)

type Store struct {
	log   *logger.Logger
	store property.Storer
	cache map[string]property.Property
	mu    sync.RWMutex
}

func NewStore(log *logger.Logger, store property.Storer) *Store {
	return &Store{
		log:   log,
		store: store,
		cache: map[string]property.Property{},
	}
}

func (s *Store) Create(ctx context.Context, prprty property.Property) error {
	if err := s.store.Create(ctx, prprty); err != nil {
		return err
	}

	s.writeCache(prprty)

	return nil
}

func (s *Store) Update(ctx context.Context, prprty property.Property) error {
	if err := s.store.Update(ctx, prprty); err != nil {
		return err
	}

	s.writeCache(prprty)

	return nil
}

func (s *Store) Delete(ctx context.Context, prprty property.Property) error {
	if err := s.store.Delete(ctx, prprty); err != nil {
		return err
	}

	s.deleteCache(prprty)

	return nil
}

func (s *Store) QueryByID(ctx context.Context, prprtyID uuid.UUID) (property.Property, error) {
	cachedUsr, ok := s.readCache(prprtyID.String())
	if ok {
		return cachedUsr, nil
	}

	usr, err := s.store.QueryByID(ctx, prprtyID)
	if err != nil {
		return property.Property{}, err
	}

	s.writeCache(usr)

	return usr, nil
}

func (s *Store) QueryByManagerID(ctx context.Context, managerID uuid.UUID) ([]property.Property, error) {
	prprties, err := s.store.QueryByManagerID(ctx, managerID)
	if err != nil {
		return nil, err
	}

	return prprties, nil
}

// =============================================================================

func (s *Store) readCache(key string) (property.Property, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	prprty, exists := s.cache[key]
	if !exists {
		return property.Property{}, false
	}

	return prprty, true
}

func (s *Store) writeCache(prprty property.Property) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cache[prprty.ID.String()] = prprty
}

func (s *Store) deleteCache(prprty property.Property) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.cache, prprty.ID.String())
}
