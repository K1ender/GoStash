package store

import (
	"errors"
	"fmt"
	"hash/fnv"
	"sync"

	"github.com/k1ender/go-stash/internal/utils"
)

type shard struct {
	m  map[string]string
	rw sync.RWMutex
}

type ShardedStore struct {
	shards []*shard
}

func NewShardedStore(numShards int) *ShardedStore {
	shards := make([]*shard, numShards)
	for i := 0; i < numShards; i++ {
		shards[i] = &shard{
			m: make(map[string]string),
		}
	}
	return &ShardedStore{
		shards: shards,
	}
}

func (s *ShardedStore) getShard(key string) *shard {
	hash := fnv.New32()
	hash.Write([]byte(key))
	return s.shards[hash.Sum32()%uint32(len(s.shards))]
}

func (s *ShardedStore) Get(key string) (string, error) {
	sh := s.getShard(key)
	sh.rw.RLock()
	defer sh.rw.RUnlock()

	if value, exists := sh.m[key]; exists {
		return value, nil
	}

	return "", ErrNotFound
}

func (s *ShardedStore) Set(key string, value string) error {
	sh := s.getShard(key)
	sh.rw.Lock()
	defer sh.rw.Unlock()

	sh.m[key] = value

	return nil
}

func (s *ShardedStore) Incr(key string) (int, error) {
	sh := s.getShard(key)
	sh.rw.Lock()
	defer sh.rw.Unlock()

	value, exists := sh.m[key]
	if !exists {
		sh.m[key] = "1"
		return 1, nil
	}

	val, err := utils.FastStringToInt(value)
	if err != nil {
		return 0, errors.New("value is not an integer")
	}

	intValue := val + 1
	sh.m[key] = fmt.Sprintf("%d", intValue)
	return intValue, nil
}

func (s *ShardedStore) Decr(key string) (int, error) {
	sh := s.getShard(key)
	sh.rw.Lock()
	defer sh.rw.Unlock()

	value, exists := sh.m[key]
	if !exists {
		sh.m[key] = "-1"
		return -1, nil
	}

	val, err := utils.FastStringToInt(value)
	if err != nil {
		return 0, errors.New("value is not an integer")
	}

	intValue := val - 1
	sh.m[key] = fmt.Sprintf("%d", intValue)
	return intValue, nil
}

func (s *ShardedStore) Del(key string) error {
	sh := s.getShard(key)
	sh.rw.Lock()
	defer sh.rw.Unlock()

	if _, exists := sh.m[key]; !exists {
		return ErrNotFound
	}

	delete(sh.m, key)
	return nil
}
