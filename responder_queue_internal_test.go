package bot

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testChat struct {
	Name string
}

func (c testChat) Load(r *Robot)            {}
func (c testChat) Username() string         { return c.Name }
func (c testChat) Messages() <-chan Message { return make(chan Message) }
func (c testChat) Send(m Message) error     { return nil }
func (c testChat) Reply(m Message) error    { return nil }
func (c testChat) Topic(m Message) error    { return nil }
func (c testChat) Direct(m Message) error   { return nil }

func TestResponderQueueCallback(t *testing.T) {
	rq := newResponderQueue(1)
	msg := Message{Text: "hi"}

	rq.On(DefaultMessage, func(rs *Responder) {
		assert.Equal(t, msg, rs.Message)
	})

	rq.Emit(DefaultMessage, &Responder{Message: msg})
}

func TestResponderQueueChannel(t *testing.T) {
	rq := newResponderQueue(2)
	msg := Message{Text: "hi"}

	pipe := rq.On(Leave)
	rq.Emit(Leave, &Responder{Message: msg})

	rs := <-pipe
	assert.Equal(t, msg, rs.Message)
}

func TestResponderQueueForward(t *testing.T) {
	rs := newResponderQueue(0)
	msg1 := Message{Type: DefaultMessage, Text: "foo"}
	msg2 := Message{Type: Leave, Text: "bar"}
	var wg sync.WaitGroup
	ch := make(chan Message, 3)
	defer close(ch)

	wg.Add(2)
	ch <- msg1
	ch <- msg2

	rs.On(DefaultMessage, func(rs *Responder) {
		assert.Equal(t, msg1, rs.Message)
		wg.Done()
	})

	rs.On(Leave, func(rs *Responder) {
		assert.Equal(t, msg2, rs.Message)
		wg.Done()
	})

	// subject
	go rs.Forward(&Robot{Chat: testChat{}}, ch)
	wg.Wait()
}

func TestMessageIsForwardedToRespondListener(t *testing.T) {
	botName := "testBot"
	tests := []struct {
		name string
		text string
	}{
		{
			name: "normal space in message",
			text: "@testBot hello",
		},
		{
			name: "no break space in message",
			text: "@testBot\u00a0hello",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rs := newResponderQueue(0)
			ch := make(chan Message, 1)
			defer close(ch)
			ch <- Message{Type: DefaultMessage, Text: test.text}

			// Channel to end the test if message fails to be received.
			doneCh := make(chan struct{})
			defer close(doneCh)

			rs.On(Response, func(rs *Responder) {
				assert.Equal(t, Message{Type: DefaultMessage, Text: test.text}, rs.Message)
				doneCh <- struct{}{}
			})

			go rs.Forward(&Robot{Chat: testChat{Name: botName}}, ch)
			select {
			case <-time.NewTicker(time.Second).C:
				assert.Fail(t, "Timeout waiting to receive message")
			case <-doneCh:
			}

		})
	}
}
