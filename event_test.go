package gobot_test

import (
	"sync"
	"testing"

	"github.com/berfarah/gobot"
	"github.com/stretchr/testify/assert"
)

type eventDispatcherRun struct {
	mu    sync.Mutex
	Count int
}

func (r *eventDispatcherRun) Increment() {
	r.mu.Lock()
	r.Count++
	r.mu.Unlock()
}

func TestEventDispatcher(t *testing.T) {
	run := &eventDispatcherRun{}
	ed := gobot.NewEventDispatcher()

	ed.On("some:event", func(d interface{}) {})
	ed.On("success", func(d interface{}) { run.Increment() })

	ed.Emit("some:event", nil)
	ed.Emit("some:event", nil)
	ed.Emit("some:event", nil)

	ed.Emit("success", nil)
	ed.Emit("success", nil)

	ed.Wait()
	assert.Equal(t, 2, run.Count)
}
