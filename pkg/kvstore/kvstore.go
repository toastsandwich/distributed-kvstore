package kvstore

import (
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

type EventTrigger func(string, any) error

var (
	AddEvent    EventTrigger
	UpdateEvent EventTrigger
	RemoveEvent EventTrigger
)

type Store struct {
	rwmu sync.RWMutex
	KV   map[string]any
}

func New() *Store {
	return &Store{
		rwmu: sync.RWMutex{},
		KV:   make(map[string]any),
	}
}

func (s *Store) validateKey(k string) error {
	if k == "" {
		return ErrEmptyKey
	}

	if slices.Contains(restrictedKeys, k) {
		return ErrRestrictedAccess
	}
	return nil
}

func (s *Store) has(k string) bool {
	_, ok := s.KV[k]
	return ok
}

func (s *Store) Get(k string) (any, error) {
	s.rwmu.RLock()
	defer s.rwmu.RUnlock()

	if err := s.validateKey(k); err != nil {
		return nil, err
	}

	if !s.has(k) {
		return nil, ErrKeyDoesNotExist
	}

	return s.KV[k], nil
}

func (s *Store) Add(k string, v any, sync bool) error {
	s.rwmu.Lock()
	defer s.rwmu.Unlock()

	if err := s.validateKey(k); err != nil {
		return err
	}

	if s.has(k) {
		return ErrKeyAlreadyExist
	}

	s.KV[k] = v

	if sync {
		return AddEvent(k, v)
	}
	return nil
}

func (s *Store) Update(k string, v any, sync bool) error {
	s.rwmu.Lock()
	defer s.rwmu.Unlock()

	if err := s.validateKey(k); err != nil {
		return err
	}

	if !s.has(k) {
		return ErrKeyDoesNotExist
	}

	s.KV[k] = v
	if sync {
		return UpdateEvent(k, v)
	}
	return nil
}

func (s *Store) Remove(k string, sync bool) error {
	s.rwmu.Lock()
	defer s.rwmu.Unlock()

	if err := s.validateKey(k); err != nil {
		return err
	}

	if !s.has(k) {
		return ErrKeyDoesNotExist
	}

	delete(s.KV, k)
	if sync {
		return RemoveEvent(k, "")
	}
	return nil
}
