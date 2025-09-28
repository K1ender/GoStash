package store

import (
	"errors"
	"fmt"
	"sync"

	"github.com/k1ender/go-stash/internal/utils"
)

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

func (s *HashMapStore) Incr(key string) (int, error) {
	s.rw.Lock()
	defer s.rw.Unlock()
	value, exists := s.data[key]
	if !exists {
		s.data[key] = "1"
		return 1, nil
	}

	val, err := utils.FastStringToInt(value)
	if err != nil {
		return 0, errors.New("value is not an integer")
	}

	intValue := val + 1
	s.data[key] = fmt.Sprintf("%d", intValue)
	return intValue, nil
}

func (s *HashMapStore) Decr(key string) (int, error) {
	s.rw.Lock()
	defer s.rw.Unlock()

	value, exists := s.data[key]
	if !exists {
		s.data[key] = "-1"
		return -1, nil
	}

	val, err := utils.FastStringToInt(value)
	if err != nil {
		return 0, errors.New("value is not an integer")
	}

	intValue := val - 1
	s.data[key] = fmt.Sprintf("%d", intValue)
	return intValue, nil
}
