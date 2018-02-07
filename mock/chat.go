package mock

import "github.com/botopolis/bot"

// Chat is a mock bot.Chat
type Chat struct {
	*Plugin
	// Returned by Username
	Name string
	// Returned by Messages
	MessageChan chan bot.Message
	// Called by Send
	SendFunc func(bot.Message) error
	// Called by Reply
	ReplyFunc func(bot.Message) error
	// Called by Direct
	DirectFunc func(bot.Message) error
	// Called by Topic
	TopicFunc func(bot.Message) error
}

// NewChat returns a NOOP Chat
func NewChat() *Chat {
	return &Chat{
		Plugin:     NewPlugin(),
		SendFunc:   func(bot.Message) error { return nil },
		ReplyFunc:  func(bot.Message) error { return nil },
		DirectFunc: func(bot.Message) error { return nil },
		TopicFunc:  func(bot.Message) error { return nil },
	}
}

// Username returns Chat.Name
func (a *Chat) Username() string { return a.Name }

// Messages returns Chat.MessageChan
func (a *Chat) Messages() <-chan bot.Message { return a.MessageChan }

// Send delegates to Chat.SendFunc
func (a *Chat) Send(m bot.Message) error { return a.SendFunc(m) }

// Reply delegates to Chat.ReplyFunc
func (a *Chat) Reply(m bot.Message) error { return a.ReplyFunc(m) }

// Direct delegates to Chat.DirectFunc
func (a *Chat) Direct(m bot.Message) error { return a.DirectFunc(m) }

// Topic delegates to Chat.TopicFunc
func (a *Chat) Topic(m bot.Message) error { return a.TopicFunc(m) }
