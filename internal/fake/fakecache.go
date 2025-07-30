package fake

import (
	"math/rand"
	"sync"
	"time"
)

type FakeCache struct {
	lock  sync.RWMutex
	rand  *rand.Rand
	cache map[string][]byte
}

func New() *FakeCache {
	return &FakeCache{
		rand:  rand.New(rand.NewSource(time.Now().UnixNano())),
		cache: make(map[string][]byte),
	}
}

// Get retrieves a value by key.
func (fc *FakeCache) Get(key string) ([]byte, error) {
	fc.lock.RLock() // read lock
	if item, ok := fc.cache[key]; ok {
		fc.lock.RUnlock()
		return item, nil
	}
	fc.lock.RUnlock()
	return nil, nil
}

// Set stores a value by key.
func (fc *FakeCache) Set(key string, value []byte) error {
	fc.lock.Lock() // write lock
	fc.cache[key] = value
	fc.lock.Unlock()
	return nil
}

// Delete removes a key.
func (fc *FakeCache) Delete(key string) error {
	fc.lock.Lock()
	delete(fc.cache, key)
	fc.lock.Unlock()
	return nil
}
