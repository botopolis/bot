package gobot

import "github.com/nlopes/slack"

// Responder is the object one receives when listening for an event
type Responder struct {
	Robot *Robot

	Message *slack.MessageEvent

	// Parts of Message
	UserID  string
	Channel string
	Text    string

	// Only relevant for Regexp
	Match []string
}

type Message struct {
	Channel string
	Text    string
	Params  slack.PostMessageParameters
}

func newResponder(r *Robot, m *slack.MessageEvent) *Responder {
	return &Responder{
		Robot:   r,
		Message: m,
		Channel: m.Channel,
		Text:    m.Text,
		Match:   []string{m.Text},
	}
}

func (r *Responder) User() (*User, error) {
	return &User{}, nil
}

func (r *Responder) Send(m Message) error {
	if m.Channel == "" {
		m.Channel = r.Channel
	}
	return r.Robot.Send(m)
}

func (r *Responder) Reply(m Message) error {
	m.Text = "<@" + r.UserID + "> " + m.Text
	return r.Send(m)
}

func (r *Responder) React(reaction string) error {
	return r.Robot.React(r.Channel, r.Message.Timestamp, reaction)
}

func (r *Responder) Topic(topic string) error {
	return r.Robot.SetTopic(r.Channel, topic)
}
