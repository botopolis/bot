package bot_test

import (
	"testing"

	"github.com/botopolis/bot"
	"github.com/botopolis/bot/mock"
	"github.com/stretchr/testify/assert"
)

func TestResponderSend(t *testing.T) {
	envelope := bot.Message{Room: "general"}
	message := bot.Message{Text: "hi", Room: envelope.Room, Envelope: envelope}
	chat := mock.NewChat()
	robot := bot.Robot{Chat: chat}
	res := bot.Responder{Robot: &robot, Message: message}

	chat.SendFunc = func(m bot.Message) error {
		assert.Equal(t, envelope, m.Envelope, "should receive the previous message as an envelope")
		assert.Equal(t, "hi", m.Text, "should receive the original parameters")
		assert.Equal(t, "general", m.Room, "should default the room when empty")
		return nil
	}
	res.Send(bot.Message{Text: "hi"})
}

func TestResponderReply(t *testing.T) {
	envelope := bot.Message{Room: "general"}
	message := bot.Message{Text: "hi", Envelope: envelope}
	chat := mock.NewChat()
	robot := bot.Robot{Chat: chat}
	res := bot.Responder{Robot: &robot, Message: message}

	chat.ReplyFunc = func(m bot.Message) error {
		assert.Equal(t, envelope, m.Envelope, "should receive the previous message as an envelope")
		assert.Equal(t, "hi", m.Text, "should receive the original parameters")
		return nil
	}
	res.Reply("hi")
}

func TestResponderTopic(t *testing.T) {
	envelope := bot.Message{Room: "general"}
	message := bot.Message{Text: "hi", Envelope: envelope}
	chat := mock.NewChat()
	robot := bot.Robot{Chat: chat}
	res := bot.Responder{Robot: &robot, Message: message}

	chat.TopicFunc = func(m bot.Message) error {
		assert.Equal(t, envelope, m.Envelope, "should receive the previous message as an envelope")
		assert.Equal(t, "stuff", m.Topic, "should receive the original parameters")
		return nil
	}
	res.Topic("stuff")
}

func TestUser(t *testing.T) {
	name := "bob"
	not := "notbob"
	res := bot.Responder{Message: bot.Message{User: name}}

	assert.False(t, bot.User(not)(&res), "is false for differing user")
	assert.True(t, bot.User(name)(&res), "is true for matching user")
}

func TestRoom(t *testing.T) {
	name := "general"
	not := "notgeneral"
	res := bot.Responder{Message: bot.Message{Room: name}}

	assert.False(t, bot.Room(not)(&res), "is false for differing room")
	assert.True(t, bot.Room(name)(&res), "is true for matching room")
}

func TestContains(t *testing.T) {
	text := "hi there bob"
	yes := "bob"
	not := "notbob"
	res := bot.Responder{Message: bot.Message{Text: text}}

	assert.False(t, bot.Contains(not)(&res), "is false for excluded text")
	assert.True(t, bot.Contains(yes)(&res), "is true for contained text")
}

func TestRegexp(t *testing.T) {
	text := "hi there bob"
	yes := "b([\\w])b"
	not := "asdf(g)"
	res := bot.Responder{Message: bot.Message{Text: text}}

	assert.False(t, bot.Regexp(not)(&res), "is false for unmatched text")
	assert.Empty(t, res.Match, "doesn't create capture groups when unmatched")

	assert.True(t, bot.Regexp(yes)(&res), "is true for matched text")
	assert.Equal(t, []string{"bob", "o"}, res.Match, "creates capture groups when matched")
}
