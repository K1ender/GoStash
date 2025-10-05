package store

import (
	"errors"
	"runtime"
	"strconv"
	"sync"

	"github.com/k1ender/go-stash/internal/utils"
)

type shard struct {
	m  map[string]string
	rw sync.RWMutex
}

type ShardedStore struct {
	shards    []*shard
	numShards int
}

func NewShardedStore(numShards int) *ShardedStore {
	if numShards <= 0 {
		numShards = runtime.GOMAXPROCS(0) * 4
	}

	shards := make([]*shard, numShards)
	for i := range numShards {
		shards[i] = &shard{
			m: make(map[string]string),
		}
	}
	return &ShardedStore{
		shards:    shards,
		numShards: numShards,
	}
}

func fastHash(s string) uint32 {
	h := uint32(2166136261)
	for i := 0; i < len(s); i++ {
		h ^= uint32(s[i])
		h *= 16777619
	}
	return h
}

func (s *ShardedStore) getShard(key string) *shard {
	hash := fastHash(key)
	return s.shards[hash%uint32(s.numShards)]
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
	sh.m[key] = value
	sh.rw.Unlock()

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
	sh.m[key] = strconv.Itoa(intValue)
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
	sh.m[key] = strconv.Itoa(intValue)
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
