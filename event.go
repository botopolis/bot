package gobot

import "sync"

// On binds a listener to an event
func (r *Robot) On(on string, fn func(interface{})) {
	r.events.On(on, fn)
}

// Emit allows for dispatching an event for consumption by listeners
func (r *Robot) Emit(on string, ev interface{}) {
	r.events.Emit(on, ev)
}

// EventDispatcher is responsible for holding listeners and
// dispatching events to them asynchronously
type EventDispatcher struct {
	mu    sync.Mutex
	wg    sync.WaitGroup
	hooks map[string][]func(interface{})
}

// NewEventDispatcher creates an EventDispatcher
func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{hooks: make(map[string][]func(interface{}))}
}

// On binds a listener to the EventDispatcher
func (e *EventDispatcher) On(on string, fn func(interface{})) {
	e.mu.Lock()
	e.hooks[on] = append(e.hooks[on], fn)
	e.mu.Unlock()
}

func (e *EventDispatcher) exec(ev interface{}, fn func(interface{})) {
	e.wg.Add(1)
	go func() {
		fn(ev)
		e.wg.Done()
	}()
}

// Emit dispatches an event to corresponding listeners
func (e *EventDispatcher) Emit(on string, ev interface{}) {
	for _, fn := range e.hooks[on] {
		e.exec(ev, fn)
	}
}

// Wait blocks until current listeners have completed operating
func (e *EventDispatcher) Wait() {
	e.wg.Wait()
}
