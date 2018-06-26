package bot_test

import (
	"testing"

	"github.com/botopolis/bot"
	"github.com/botopolis/bot/mock"
	"github.com/stretchr/testify/assert"
)

func TestBrain(t *testing.T) {
	assert := assert.New(t)
	brain := bot.NewBrain()

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
	store := mock.NewStore()
	brain := bot.NewBrain()
	brain.SetStore(store)

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
	store := mock.NewStore()
	brain := bot.NewBrain()
	brain.SetStore(store)

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
	store := mock.NewStore()
	brain := bot.NewBrain()
	brain.SetStore(store)

	var called bool
	store.DeleteFunc = func(key string) error {
		called = true
		return nil
	}

	brain.Delete("key")
	assert.True(called, "should also delete in the store")
}

func TestBrainKeys(t *testing.T) {
	assert := assert.New(t)
	brain := bot.NewBrain()

	var noKeys []string
	assert.Equal(noKeys, brain.Keys(), "should not have any keys")

	brain.Set("A. key", nil)
	assert.Equal([]string{"A. key"}, brain.Keys(), "should have one key")

	brain.Set("C. another", nil)
	assert.Equal([]string{"A. key", "C. another"}, brain.Keys(), "should have two keys")

	brain.Set("B. alphabetized", nil)
	assert.Equal([]string{"A. key", "B. alphabetized", "C. another"}, brain.Keys(), "should alphabetize keys")
}
