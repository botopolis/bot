package gobot

import "regexp"

// Chat is the interface for chat integrations
type Chat interface {
	Plugin
	Username() string
	Messages() <-chan Message
	Send(Message) error
	Reply(Message) error
	Topic(Message) error
}

const (
	MessageEvent EventType = "chat:message"
	RespondEvent EventType = "chat:respond"
	EnterEvent   EventType = "chat:enter"
	LeaveEvent   EventType = "chat:leave"
	TopicEvent   EventType = "chat:topic"
)

type Message struct {
	Event EventType

	User  string
	Room  string
	Text  string
	Topic string

	// The original message (ptr)
	Envelope interface{}
	// Extra information
	Extra interface{}
}

// Responder is the object one receives when listening for an event
type Responder struct {
	*Robot
	Message

	// Only relevant for Regexp
	Match []string
}

func newResponder(r *Robot, m Message) *Responder {
	return &Responder{
		Robot:   r,
		Message: m,
		Match:   []string{m.Text},
	}
}

func (r *Robot) Username() string { return r.Chat.Username() }

func (r *Responder) Send(m Message) error {
	if m.Room == "" {
		m.Room = r.Room
	}
	return r.Chat.Send(m)
}

func (r *Responder) Reply(m Message) error {
	return r.Chat.Reply(m)
}

func (r *Responder) Topic(topic string) error {
	msg := r.Message
	msg.Topic = topic
	return r.Chat.Topic(msg)
}

type Hook struct {
	Name  string
	Func  func(r Responder) error
	Match Matcher
}

type Matcher func(r *Responder) bool

func (h *Hook) Run(r *Responder) error {
	if h.Match != nil && !h.Match(r) {
		return nil
	}
	if err := h.Func(*r); err != nil {
		return err
	}
	return nil
}

func MatchUser(u string) Matcher {
	return func(r *Responder) bool {
		return r.Message.User == u
	}
}

func MatchRoom(c string) Matcher {
	return func(r *Responder) bool {
		return r.Message.Room == c
	}
}

func MatchText(t string) Matcher {
	return func(r *Responder) bool {
		return r.Message.Text == t
	}
}

func MatchRegexp(expression string) Matcher {
	return func(r *Responder) bool {
		reg, err := regexp.Compile(expression)
		if err == nil && reg.MatchString(r.Message.Text) {
			r.Match = reg.FindStringSubmatch(r.Message.Text)
			return true
		}
		return false
	}
}
