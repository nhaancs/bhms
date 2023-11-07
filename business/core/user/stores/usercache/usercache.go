// Package usercache contains user related CRUD functionality with caching.
package usercache

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/nhaancs/bhms/business/core/user"
	"github.com/nhaancs/bhms/foundation/logger"
)

// Store manages the set of APIs for user data and caching.
type Store struct {
	log   *logger.Logger
	store user.Storer
	cache map[string]user.User
	mu    sync.RWMutex
}

// NewStore constructs the api for data and caching access.
func NewStore(log *logger.Logger, store user.Storer) *Store {
	return &Store{
		log:   log,
		store: store,
		cache: map[string]user.User{},
	}
}

// Create inserts a new user into the database.
func (s *Store) Create(ctx context.Context, core user.User) error {
	if err := s.store.Create(ctx, core); err != nil {
		return err
	}

	s.writeCache(core)

	return nil
}

// Update replaces a user document in the database.
func (s *Store) Update(ctx context.Context, core user.User) error {
	if err := s.store.Update(ctx, core); err != nil {
		return err
	}

	s.writeCache(core)

	return nil
}

// Delete removes a user from the database.
func (s *Store) Delete(ctx context.Context, core user.User) error {
	if err := s.store.Delete(ctx, core); err != nil {
		return err
	}

	s.deleteCache(core)

	return nil
}

// QueryByID gets the specified user from the database.
func (s *Store) QueryByID(ctx context.Context, id uuid.UUID) (user.User, error) {
	cachedUsr, ok := s.readCache(id.String())
	if ok {
		return cachedUsr, nil
	}

	usr, err := s.store.QueryByID(ctx, id)
	if err != nil {
		return user.User{}, err
	}

	s.writeCache(usr)

	return usr, nil
}

// QueryByIDs gets the specified users from the database.
func (s *Store) QueryByIDs(ctx context.Context, ids []uuid.UUID) ([]user.User, error) {
	usr, err := s.store.QueryByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

// QueryByPhone gets the specified user from the database by email.
func (s *Store) QueryByPhone(ctx context.Context, phone string) (user.User, error) {
	cachedUsr, ok := s.readCache(phone)
	if ok {
		return cachedUsr, nil
	}

	usr, err := s.store.QueryByPhone(ctx, phone)
	if err != nil {
		return user.User{}, err
	}
	s.writeCache(usr)

	return usr, nil
}

// =============================================================================

// readCache performs a safe search in the cache for the specified key.
func (s *Store) readCache(key string) (user.User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	usr, exists := s.cache[key]
	if !exists {
		return user.User{}, false
	}

	return usr, true
}

// writeCache performs a safe write to the cache for the specified user.
func (s *Store) writeCache(usr user.User) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cache[usr.ID.String()] = usr
	s.cache[usr.Phone] = usr
}

// deleteCache performs a safe removal from the cache for the specified user.
func (s *Store) deleteCache(usr user.User) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.cache, usr.ID.String())
	delete(s.cache, usr.Phone)
}
