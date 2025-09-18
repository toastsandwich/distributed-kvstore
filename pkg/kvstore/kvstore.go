package kvstore

import (
	"context"
	"errors"
	"slices"
	"sync"
)

var (
	ErrEmptyKey         = errors.New("key cannot be empty")
	ErrKeyDoesNotExist  = errors.New("key does not exist")
	ErrKeyAlreadyExist  = errors.New("key already exist")
	ErrRestrictedAccess = errors.New("cannot access this key")

	restrictedKeys = []string{"node"}
)

type EventTrigger func(context.Context, string, any) error

var (
	AddEvent    EventTrigger
	UpdateEvent EventTrigger
	RemoveEvent EventTrigger
)

type Store struct {
	ctx context.Context

	rwmu sync.RWMutex
	KV   map[string]any
}

func New(ctx context.Context) *Store {
	return &Store{
		ctx:  ctx,
		rwmu: sync.RWMutex{},
		KV:   make(map[string]any),
	}
}

func (s *Store) Has(k string) error {
	if k == "" {
		return ErrEmptyKey
	}

	if slices.Contains(restrictedKeys, k) {
		return ErrRestrictedAccess
	}

	if _, ok := s.KV[k]; !ok {
		return ErrKeyDoesNotExist
	}
	return nil
}

func (s *Store) Get(k string) (any, error) {
	s.rwmu.RLock()
	defer s.rwmu.RUnlock()

	if err := s.Has(k); err != nil {
		return nil, err
	}

	return s.KV[k], nil
}

func (s *Store) Add(k string, v any) error {
	s.rwmu.Lock()
	defer s.rwmu.Unlock()

	if err := s.Has(k); err != nil && err != ErrKeyDoesNotExist {
		return err
	}

	s.KV[k] = v

	return AddEvent(s.ctx, k, v)
}

func (s *Store) Update(k string, v any) error {
	s.rwmu.Lock()
	defer s.rwmu.Unlock()

	if err := s.Has(k); err != nil {
		return err
	}

	s.KV[k] = v
	return UpdateEvent(s.ctx, k, v)
}

func (s *Store) Remove(k string) error {
	s.rwmu.Lock()
	defer s.rwmu.Unlock()

	if err := s.Has(k); err != nil {
		return err
	}

	delete(s.KV, k)
	return RemoveEvent(s.ctx, k, "")
}
