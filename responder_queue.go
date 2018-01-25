package gobot

import (
	"regexp"
	"sync"
)

func newListener(capacity int, cbs ...func(*Responder)) listener {
	var cb func(*Responder)
	var ch chan Responder
	if len(cbs) > 0 {
		cb = cbs[0]
	} else {
		ch = make(chan Responder, capacity)
	}
	return listener{ch: ch, callback: cb}
}

// Listener holds the listening channel or callback
type listener struct {
	ch       chan Responder
	callback func(*Responder)

	IsClosed bool
	Once     bool
}

func (l *listener) Close() {
	l.IsClosed = true
	close(l.ch)
}

func (l *listener) Dispatch(e *Responder) {
	if l.callback != nil {
		l.callback(e)
		return
	}
	if !l.IsClosed {
		l.ch <- *e
	}
}

// responderQueue is responsible for holding listeners and
// dispatching events to them asynchronously
type responderQueue struct {
	Capacity int

	mu        sync.RWMutex
	listeners map[messageType][]listener
}

func newResponderQueue(capacity int) *responderQueue {
	return &responderQueue{
		Capacity:  capacity,
		listeners: make(map[messageType][]listener),
	}
}

func (e *responderQueue) Forward(r *Robot, ch <-chan Message) {
	for msg := range ch {
		rs := newResponder(r, msg)
		exp := regexp.MustCompile("^" + r.Username() + "\\s")
		if msg.Type == DefaultMessage && exp.MatchString(msg.Text) {
			e.Emit(Response, rs)
		}

		e.Emit(msg.Type, rs)
	}
}

// On binds a listener to the responderQueue
func (e *responderQueue) On(on messageType, cbs ...func(*Responder)) <-chan Responder {
	e.mu.Lock()
	defer e.mu.Unlock()
	l := newListener(e.Capacity, cbs...)
	e.listeners[on] = append(e.listeners[on], l)
	return l.ch
}

// Emit dispatches an event to corresponding listeners
func (e *responderQueue) Emit(on messageType, rs *Responder) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if ls, ok := e.listeners[on]; ok {
		for _, l := range ls {
			l.Dispatch(rs)
		}
	}
}
