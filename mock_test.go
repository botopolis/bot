package gobot_test

import "github.com/berfarah/gobot"

type Plugin struct {
	LoadFunc   func(r *gobot.Robot)
	UnloadFunc func(r *gobot.Robot)
}

func NewPlugin() *Plugin {
	return &Plugin{
		LoadFunc:   func(r *gobot.Robot) {},
		UnloadFunc: func(r *gobot.Robot) {},
	}
}
func (p *Plugin) Load(r *gobot.Robot)   { p.LoadFunc(r) }
func (p *Plugin) Unload(r *gobot.Robot) { p.UnloadFunc(r) }

type Chat struct {
	*Plugin
	Name        string
	MessageChan chan gobot.Message
	SendFunc    func(gobot.Message) error
	ReplyFunc   func(gobot.Message) error
	TopicFunc   func(gobot.Message) error
}

func NewChat() *Chat {
	return &Chat{
		Plugin:    NewPlugin(),
		SendFunc:  func(gobot.Message) error { return nil },
		ReplyFunc: func(gobot.Message) error { return nil },
		TopicFunc: func(gobot.Message) error { return nil },
	}
}
func (a *Chat) Username() string               { return a.Name }
func (a *Chat) Messages() <-chan gobot.Message { return a.MessageChan }
func (a *Chat) Send(m gobot.Message) error     { return a.SendFunc(m) }
func (a *Chat) Reply(m gobot.Message) error    { return a.ReplyFunc(m) }
func (a *Chat) Topic(m gobot.Message) error    { return a.TopicFunc(m) }
