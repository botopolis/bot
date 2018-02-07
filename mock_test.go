package bot_test

import "github.com/botopolis/bot"

type Plugin struct {
	LoadFunc   func(r *bot.Robot)
	UnloadFunc func(r *bot.Robot)
}

func NewPlugin() *Plugin {
	return &Plugin{
		LoadFunc:   func(r *bot.Robot) {},
		UnloadFunc: func(r *bot.Robot) {},
	}
}
func (p *Plugin) Load(r *bot.Robot)   { p.LoadFunc(r) }
func (p *Plugin) Unload(r *bot.Robot) { p.UnloadFunc(r) }

type Chat struct {
	*Plugin
	Name        string
	MessageChan chan bot.Message
	SendFunc    func(bot.Message) error
	ReplyFunc   func(bot.Message) error
	DirectFunc  func(bot.Message) error
	TopicFunc   func(bot.Message) error
}

func NewChat() *Chat {
	return &Chat{
		Plugin:     NewPlugin(),
		SendFunc:   func(bot.Message) error { return nil },
		ReplyFunc:  func(bot.Message) error { return nil },
		DirectFunc: func(bot.Message) error { return nil },
		TopicFunc:  func(bot.Message) error { return nil },
	}
}
func (a *Chat) Username() string             { return a.Name }
func (a *Chat) Messages() <-chan bot.Message { return a.MessageChan }
func (a *Chat) Send(m bot.Message) error     { return a.SendFunc(m) }
func (a *Chat) Reply(m bot.Message) error    { return a.ReplyFunc(m) }
func (a *Chat) Direct(m bot.Message) error   { return a.DirectFunc(m) }
func (a *Chat) Topic(m bot.Message) error    { return a.TopicFunc(m) }
