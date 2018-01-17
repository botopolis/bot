package memory

import (
	"errors"
	"reflect"
)

type Store struct{ Cache map[string]interface{} }

func (s *Store) Set(key string, i interface{}) error {
	s.Cache[key] = i
	return nil
}

func (s *Store) Get(key string, i interface{}) error {
	if tmp, ok := s.Cache[key]; ok {
		reflect.ValueOf(i).Elem().Set(reflect.ValueOf(tmp))
		return nil
	}
	return errors.New("Cache miss")
}

func New() *Store {
	return &Store{Cache: make(map[string]interface{})}
}
