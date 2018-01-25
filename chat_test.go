package gobot_test

import (
	"testing"

	"github.com/berfarah/gobot"
	"github.com/stretchr/testify/assert"
)

func TestResponderSend(t *testing.T) {
	envelope := gobot.Message{Room: "general"}
	chat := NewChat()
	robot := gobot.Robot{Chat: chat}
	res := gobot.Responder{Robot: &robot, Message: envelope}

	chat.SendFunc = func(m gobot.Message) error {
		assert.Equal(t, envelope, m.Envelope, "should receive the previous message as an envelope")
		assert.Equal(t, "hi", m.Text, "should receive the original parameters")
		assert.Equal(t, "general", m.Room, "should default the room when empty")
		return nil
	}
	res.Send(gobot.Message{Text: "hi"})
}

func TestResponderReply(t *testing.T) {
	envelope := gobot.Message{Room: "general"}
	chat := NewChat()
	robot := gobot.Robot{Chat: chat}
	res := gobot.Responder{Robot: &robot, Message: envelope}

	chat.ReplyFunc = func(m gobot.Message) error {
		assert.Equal(t, envelope, m.Envelope, "should receive the previous message as an envelope")
		assert.Equal(t, "hi", m.Text, "should receive the original parameters")
		return nil
	}
	res.Reply("hi")
}

func TestResponderTopic(t *testing.T) {
	envelope := gobot.Message{Room: "general"}
	chat := NewChat()
	robot := gobot.Robot{Chat: chat}
	res := gobot.Responder{Robot: &robot, Message: envelope}

	chat.TopicFunc = func(m gobot.Message) error {
		assert.Equal(t, envelope, m.Envelope, "should receive the previous message as an envelope")
		assert.Equal(t, "stuff", m.Topic, "should receive the original parameters")
		return nil
	}
	res.Topic("stuff")
}

func TestUser(t *testing.T) {
	name := "bob"
	not := "notbob"
	res := gobot.Responder{Message: gobot.Message{User: name}}

	assert.False(t, gobot.User(not)(&res), "is false for differing user")
	assert.True(t, gobot.User(name)(&res), "is true for matching user")
}

func TestRoom(t *testing.T) {
	name := "general"
	not := "notgeneral"
	res := gobot.Responder{Message: gobot.Message{Room: name}}

	assert.False(t, gobot.Room(not)(&res), "is false for differing room")
	assert.True(t, gobot.Room(name)(&res), "is true for matching room")
}

func TestContains(t *testing.T) {
	text := "hi there bob"
	yes := "bob"
	not := "notbob"
	res := gobot.Responder{Message: gobot.Message{Text: text}}

	assert.False(t, gobot.Contains(not)(&res), "is false for excluded text")
	assert.True(t, gobot.Contains(yes)(&res), "is true for contained text")
}

func TestRegexp(t *testing.T) {
	text := "hi there bob"
	yes := "b([\\w])b"
	not := "asdf(g)"
	res := gobot.Responder{Message: gobot.Message{Text: text}}

	assert.False(t, gobot.Regexp(not)(&res), "is false for unmatched text")
	assert.Empty(t, res.Match, "doesn't create capture groups when unmatched")

	assert.True(t, gobot.Regexp(yes)(&res), "is true for matched text")
	assert.Equal(t, []string{"bob", "o"}, res.Match, "creates capture groups when matched")
}
