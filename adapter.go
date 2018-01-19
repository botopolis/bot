package gobot

type MessageEvent struct {
	Type string
	Data interface{}
}

type Message struct {
	User string
	Room string
	Text string

	Extra interface{}
}

type EnterMessage Message
type LeaveMessage Message
type TopicChange struct {
	Message
	Topic string
}

type Adapter interface {
	Messages() <-chan *MessageEvent
	Send(Message) error
	Reply(Message) error
	Topic(TopicChange) error
}
