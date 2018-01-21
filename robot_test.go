package gobot_test

import (
	"testing"

	"github.com/berfarah/gobot"
	"github.com/stretchr/testify/assert"
)

func TestRobotNew(t *testing.T) {
	var loaded bool
	chat := NewChat()
	chat.Plugin.LoadFunc = func(r *gobot.Robot) { loaded = true }

	gobot.New(chat)
	assert.True(t, loaded, "Calls the loaded function on the adapter")
}
