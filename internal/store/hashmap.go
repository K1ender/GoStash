package store

import "sync"

type HashMapStore struct {
	data map[string]string
	rw   sync.RWMutex
}

func NewHashMapStore() *HashMapStore {
	return &HashMapStore{
		data: make(map[string]string),
	}
}

func (s *HashMapStore) Get(key string) (string, error) {
	s.rw.RLock()
	defer s.rw.RUnlock()
	if value, exists := s.data[key]; exists {
		return value, nil
	}
	return "", ErrNotFound
}

func (s *HashMapStore) Set(key, value string) error {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.data[key] = value
	return nil
}
