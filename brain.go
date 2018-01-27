package gobot

import (
	"fmt"
	"reflect"
	"sync"
)

// Store is what the brain saves information in for longer term
// recollection. Brain loads what it needs to in memory, but in
// order to have persistence between runs, we can use a store.
type Store interface {
	// Get is called on Brain.Get() when the key doesn't exist in memory
	Get(key string, object interface{}) error
	// Set is always called on Brain.Set()
	Set(key string, object interface{}) error
	// Delete is always called on Brain.Delete()
	Delete(key string) error
}

// Brain is our data store. It defaults to just saving values in memory,
// but can be given a Store via Brain.SetStore(), to which values are written.
// Brain is threadsafe - and assumes the same of Store.
type Brain struct {
	store Store
	mu    sync.RWMutex
	cache map[string]interface{}
}

type nullStore struct{}
type nullStoreError struct{}

func (e nullStoreError) Error() string { return "Nullstore contains no values" }

func (n nullStore) Get(key string, object interface{}) error { return nullStoreError{} }
func (n nullStore) Set(key string, object interface{}) error { return nullStoreError{} }
func (n nullStore) Delete(key string) error                  { return nullStoreError{} }

// NewBrain creates a Brain
func NewBrain() *Brain {
	return &Brain{
		store: nullStore{},
		cache: make(map[string]interface{}),
	}
}

// SetStore assigns a Store to Brain
func (b *Brain) SetStore(s Store) { b.store = s }

// Get retrieves a value from the store and sets it to the interface
// It tries the memory store first, and falls back to Brain's Store
func (b *Brain) Get(key string, i interface{}) error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if d, ok := b.cache[key]; ok {
		return copyInterface(d, i)
	}
	if err := b.store.Get(key, i); err != nil {
		return err
	}
	b.cache[key] = i
	return nil
}

// Set assigns the given value in memory and to Brain's Store
func (b *Brain) Set(key string, i interface{}) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.cache[key] = i
	return b.store.Set(key, i)
}

// Delete removes the given key from memory and Brain's Store
func (b *Brain) Delete(key string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.cache, key)
	return b.store.Delete(key)
}

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
