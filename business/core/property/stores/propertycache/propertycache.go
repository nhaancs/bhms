package propertycache

import (
	"context"
	"github.com/google/uuid"
	"github.com/nhaancs/bhms/business/core/property"
	"github.com/nhaancs/bhms/foundation/logger"
	"sync"
)

type Store struct {
	log            *logger.Logger
	store          property.Storer
	cache          map[string]property.Property
	userProperties map[string][]string
	mu             sync.RWMutex
}

func NewStore(log *logger.Logger, store property.Storer) *Store {
	return &Store{
		log:            log,
		store:          store,
		cache:          map[string]property.Property{},
		userProperties: map[string][]string{},
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

func (s *Store) QueryByID(ctx context.Context, id uuid.UUID) (property.Property, error) {
	cached, ok := s.readCache(id.String())
	if ok {
		return cached, nil
	}

	prprty, err := s.store.QueryByID(ctx, id)
	if err != nil {
		return property.Property{}, err
	}

	s.writeCache(prprty)

	return prprty, nil
}

func (s *Store) QueryByManagerID(ctx context.Context, managerID uuid.UUID) ([]property.Property, error) {
	prprtyIDs, ok := s.readCacheByManagerID(managerID.String())
	if ok {
		var prprties []property.Property
		for _, prprtyID := range prprtyIDs {
			p, ok := s.readCache(prprtyID)
			if !ok {
				break
			}
			prprties = append(prprties, p)
		}
		if len(prprtyIDs) == len(prprties) {
			return prprties, nil
		}
	}

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

func (s *Store) readCacheByManagerID(key string) ([]string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	prprtyIDs, exists := s.userProperties[key]
	if !exists {
		return nil, false
	}

	return prprtyIDs, true
}

func (s *Store) writeCache(prprty property.Property) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cache[prprty.ID.String()] = prprty

	prprtyIDs, _ := s.userProperties[prprty.ManagerID.String()]
	for _, prprtyID := range prprtyIDs {
		if prprtyID == prprty.ID.String() {
			return
		}
	}
	s.userProperties[prprty.ManagerID.String()] = append(prprtyIDs, prprty.ID.String())
}

func (s *Store) deleteCache(prprty property.Property) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.cache, prprty.ID.String())

	if prprtyIDs, exist := s.userProperties[prprty.ManagerID.String()]; exist {
		var left []string
		for _, prprtyID := range prprtyIDs {
			if prprtyID == prprty.ID.String() {
				continue
			}
			left = append(left, prprtyID)
		}
		s.userProperties[prprty.ManagerID.String()] = left
	}
}
