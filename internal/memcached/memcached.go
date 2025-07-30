package memcached

import (
	"github.com/a1ndreay/memproxy/pkg/cache"
	"github.com/bradfitz/gomemcache/memcache"
)

type Memcached struct {
	client *memcache.Client
}

func New(addr string) cache.Backend {
	return &Memcached{client: memcache.New(addr)}
}

// Get retrieves a value by key.
func (m *Memcached) Get(key string) ([]byte, error) {
	item, err := m.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return nil, nil
		}
		return nil, err
	}
	return item.Value, nil
}

// Set stores a value by key.
func (m *Memcached) Set(key string, value []byte) error {
	return m.client.Set(&memcache.Item{Key: key, Value: value})
}

// Delete removes a key.
func (m *Memcached) Delete(key string) error {
	return m.client.Delete(key)
}
