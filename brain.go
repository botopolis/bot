package gobot

import (
	"fmt"
	"reflect"
	"sync"
)

type Store interface {
	Get(key string, object interface{}) error
	Set(key string, object interface{}) error
	Delete(key string) error
	Save() error
}

type Brain struct {
	Store Store
	mu    sync.RWMutex
	cache map[string]interface{}
}

type nullStore struct{}

func (n nullStore) Get(key string, object interface{}) error { return nil }
func (n nullStore) Set(key string, object interface{}) error { return nil }
func (n nullStore) Delete(key string) error                  { return nil }
func (n nullStore) Save() error                              { return nil }

func newBrain() *Brain {
	return &Brain{
		Store: nullStore{},
		cache: make(map[string]interface{}),
	}
}

func (b *Brain) Get(key string, i interface{}) error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if d, ok := b.cache[key]; ok {
		return copyInterface(d, i)
	}
	if err := b.Store.Get(key, i); err != nil {
		return err
	}
	b.cache[key] = i
	return nil
}

func (b *Brain) Set(key string, i interface{}) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.cache[key] = i
	return b.Store.Set(key, i)
}

func (b *Brain) Delete(key string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.cache, key)
	return b.Store.Delete(key)
}

func (b *Brain) Save() error      { return b.Store.Save() }
func (b *Brain) SetStore(s Store) { b.Store = s }

func ifaceValue(i interface{}) reflect.Value {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		return reflect.Indirect(v)
	}
	return v
}

func copyInterface(from, to interface{}) error {
	if to == nil {
		return nil
	}
	fromVal, toVal := ifaceValue(from), ifaceValue(to)
	if fromVal.Kind() != toVal.Kind() {
		return fmt.Errorf("Can't assign %T to %T", from, to)
	}
	toVal.Set(fromVal)
	return nil
}
