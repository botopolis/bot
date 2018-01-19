package gobot

// Responder is the object one receives when listening for an event
type Responder struct {
	Robot *Robot

	Message Message

	// Parts of Message
	User string
	Room string
	Text string

	// Only relevant for Regexp
	Match []string
}

func newResponder(r *Robot, m Message) *Responder {
	return &Responder{
		Robot:   r,
		Message: m,
		User:    m.User,
		Room:    m.Room,
		Text:    m.Text,
		Match:   []string{m.Text},
	}
}

func (r *Responder) Send(m Message) error {
	if m.Room == "" {
		m.Room = r.Room
	}
	return r.Robot.Adapter.Send(m)
}

func (r *Responder) Reply(m Message) error {
	return r.Robot.Adapter.Reply(m)
}

func (r *Responder) Topic(topic string) error {
	return r.Robot.Adapter.Topic(TopicChange{Message: r.Message, Topic: topic})
}
