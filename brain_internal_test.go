package gobot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStore struct {
	GetFunc    func(string, interface{}) error
	SetFunc    func(string, interface{}) error
	DeleteFunc func(string) error
	SaveFunc   func() error
}

func newTestStore() *testStore {
	return &testStore{
		GetFunc:    func(s string, i interface{}) error { return nil },
		SetFunc:    func(s string, i interface{}) error { return nil },
		DeleteFunc: func(s string) error { return nil },
		SaveFunc:   func() error { return nil },
	}
}
func (s testStore) Get(k string, i interface{}) error { return s.GetFunc(k, i) }
func (s testStore) Set(k string, i interface{}) error { return s.SetFunc(k, i) }
func (s testStore) Delete(k string) error             { return s.DeleteFunc(k) }
func (s testStore) Save() error                       { return s.SaveFunc() }

func TestBrain(t *testing.T) {
	assert := assert.New(t)
	brain := newBrain()

	key := "foo"
	value := "bar"

	var t1 string
	brain.Set(key, value)
	brain.Get(key, &t1)
	assert.Equal(value, t1)

	var t2 string
	brain.Delete(key)
	brain.Get(key, &t2)
	assert.Equal("", t2)
}

func TestBrainGet(t *testing.T) {
	assert := assert.New(t)
	store := newTestStore()
	brain := newBrain()
	brain.Store = store

	called := map[string]int{}
	store.GetFunc = func(key string, i interface{}) error {
		called[key] = called[key] + 1
		return nil
	}

	brain.Get("missing", nil)
	assert.Equal(1, called["missing"], "uses the store when the value is not cached")

	brain.Get("missing", nil)
	assert.Equal(1, called["missing"], "caches the result from the first missing Get()")

	brain.Set("present", nil)
	brain.Get("present", nil)
	assert.True(called["present"] == 0, "doesn't use the store when the value is cached")
}

func TestBrainSet(t *testing.T) {
	assert := assert.New(t)
	store := newTestStore()
	brain := newBrain()
	brain.Store = store

	var called bool
	store.SetFunc = func(key string, i interface{}) error {
		called = true
		return nil
	}

	brain.Set("key", nil)
	assert.True(called, "should also set in the store")
}

func TestBrainDelete(t *testing.T) {
	assert := assert.New(t)
	store := newTestStore()
	brain := newBrain()
	brain.Store = store

	var called bool
	store.DeleteFunc = func(key string) error {
		called = true
		return nil
	}

	brain.Delete("key")
	assert.True(called, "should also delete in the store")
}

func TestBrainSave(t *testing.T) {
	assert := assert.New(t)
	store := newTestStore()
	brain := newBrain()
	brain.Store = store

	var called bool
	store.SaveFunc = func() error {
		called = true
		return nil
	}

	brain.Save()
	assert.True(called, "should call save on store")
}
