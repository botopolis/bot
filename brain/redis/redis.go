package redis

import (
	"github.com/go-redis/cache"
	"github.com/go-redis/redis"
	"github.com/vmihailenco/msgpack"
)

type Store struct{ cache *cache.Codec }

func (s *Store) Set(key string, i interface{}) error {
	return s.cache.Set(&cache.Item{Key: key, Object: i})
}

func (s *Store) Get(key string, i interface{}) error {
	return s.cache.Get(key, i)
}

func New(address string) *Store {
	return &Store{
		cache: &cache.Codec{
			Redis: redis.NewClient(&redis.Options{
				Addr: address,
			}),
			Marshal: func(v interface{}) ([]byte, error) {
				return msgpack.Marshal(v)
			},
			Unmarshal: func(b []byte, v interface{}) error {
				return msgpack.Unmarshal(b, v)
			},
		},
	}
}
