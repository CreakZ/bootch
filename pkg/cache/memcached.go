package cache

import (
	"bootch/internal/cfg"
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
)

type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
}

type memcached struct {
	exp    int32
	client *memcache.Client
}

func NewCache(cfg cfg.Config) Cache {
	addr := fmt.Sprintf("%s:%d", cfg.MemcachedHost, cfg.MemcachedPort)
	client := memcache.New(addr)

	return memcached{
		exp:    cfg.CacheExpiration,
		client: client,
	}
}

func (m memcached) Get(key string) ([]byte, error) {
	item, err := m.client.Get(key)
	if err != nil {
		return []byte{}, err
	}
	return item.Value, nil
}

func (m memcached) Set(key string, value []byte) error {
	err := m.client.Set(&memcache.Item{
		Key:        key,
		Value:      value,
		Expiration: m.exp,
	})
	if err != nil {
		return err
	}
	return nil
}
