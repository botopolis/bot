package gobot

import (
	"regexp"
	"strings"
)

// Chat is the interface for chat integrations
type Chat interface {
	Plugin
	Username() string
	Messages() <-chan Message
	Send(Message) error
	Reply(Message) error
	Topic(Message) error
}

type messageType int
type hook func(r Responder) error

// Matcher determines whether the bot is triggered
type Matcher func(r *Responder) bool

const (
	// DefaultMessage assigns our message to robot.Hear()
	DefaultMessage messageType = iota
	// Response assigns our message to robot.Respond()
	Response
	// Enter assigns our message to robot.Enter()
	Enter
	// Leave assigns our message to robot.Leave()
	Leave
	// Topic assigns our message to robot.Topic()
	Topic
)

// Message is our message wrapper
type Message struct {
	Type  messageType
	User  string
	Room  string
	Text  string
	Topic string

	// Envelope is the original message that the bot is reacting to
	Envelope interface{}
	// Params is for adapter-specific parameters
	Params interface{}
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

// Send sends a message. It defaults to sending in the current room
func (r *Responder) Send(m Message) error {
	m.Envelope = r.Message
	if m.Room == "" {
		m.Room = r.Room
	}
	return r.Chat.Send(m)
}

// Reply responds to a message
func (r *Responder) Reply(text string) error {
	return r.Chat.Reply(Message{
		Text:     text,
		Room:     r.Room,
		Envelope: r.Message,
	})
}

// Topic changes the topic
func (r *Responder) Topic(topic string) error {
	return r.Chat.Topic(Message{
		Room:     r.Room,
		Topic:    topic,
		Envelope: r.Message,
	})
}

// User is a Matcher function that matches the User exactly
func User(u string) Matcher {
	return func(r *Responder) bool {
		return r.Message.User == u
	}
}

// Room is a Matcher function that matches the Room exactly
func Room(c string) Matcher {
	return func(r *Responder) bool {
		return r.Message.Room == c
	}
}

// Contains matches a subset of the message text
func Contains(t string) Matcher {
	return func(r *Responder) bool {
		return strings.Contains(r.Message.Text, t)
	}
}

// Regexp matches the message text with a regular expression
func Regexp(expression string) Matcher {
	return func(r *Responder) bool {
		reg, err := regexp.Compile(expression)
		if err == nil && reg.MatchString(r.Message.Text) {
			r.Match = reg.FindStringSubmatch(r.Message.Text)
			return true
		}
		return false
	}
}
