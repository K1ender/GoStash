package store

import "errors"

var (
	ErrNotFound = errors.New("key not found")
)

type Store interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Incr(key string) (int, error)
	Decr(key string) (int, error)
	Del(key string) error
}
