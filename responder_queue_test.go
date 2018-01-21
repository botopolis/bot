package gobot_test

import (
	"sync"
	"testing"

	"github.com/berfarah/gobot"
	"github.com/stretchr/testify/assert"
)

func TestResponderQueueCallback(t *testing.T) {
	rq := gobot.NewResponderQueue(1)
	msg := gobot.Message{Text: "hi"}

	rq.On(gobot.MessageEvent, func(rs *gobot.Responder) {
		assert.Equal(t, msg, rs.Message)
	})

	rq.Emit(gobot.MessageEvent, &gobot.Responder{Message: msg})
}

func TestResponderQueueChannel(t *testing.T) {
	rq := gobot.NewResponderQueue(2)
	msg := gobot.Message{Text: "hi"}

	pipe := rq.On(gobot.LeaveEvent)
	rq.Emit(gobot.LeaveEvent, &gobot.Responder{Message: msg})

	rs := <-pipe
	assert.Equal(t, msg, rs.Message)
}

func TestResponderQueueForward(t *testing.T) {
	rs := gobot.NewResponderQueue(0)
	msg1 := gobot.Message{Event: gobot.MessageEvent, Text: "foo"}
	msg2 := gobot.Message{Event: gobot.LeaveEvent, Text: "bar"}
	var wg sync.WaitGroup
	ch := make(chan gobot.Message, 3)
	defer close(ch)

	wg.Add(2)
	ch <- msg1
	ch <- msg2

	rs.On(gobot.MessageEvent, func(rs *gobot.Responder) {
		assert.Equal(t, msg1, rs.Message)
		wg.Done()
	})

	rs.On(gobot.LeaveEvent, func(rs *gobot.Responder) {
		assert.Equal(t, msg2, rs.Message)
		wg.Done()
	})

	// subject
	go rs.Forward(&gobot.Robot{Chat: NewChat()}, ch)
	wg.Wait()
}
