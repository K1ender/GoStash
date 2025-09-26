package store

type HashMapStore struct {
	data map[string]string
}

func NewHashMapStore() *HashMapStore {
	return &HashMapStore{
		data: make(map[string]string),
	}
}

func (s *HashMapStore) Get(key string) (string, error) {
	if value, exists := s.data[key]; exists {
		return value, nil
	}
	return "", ErrNotFound
}

func (s *HashMapStore) Set(key, value string) error {
	s.data[key] = value
	return nil
}
