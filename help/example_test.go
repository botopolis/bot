package help_test

import (
	"fmt"

	"github.com/botopolis/bot"
	"github.com/botopolis/bot/help"
	"github.com/botopolis/bot/mock"
)

type HelperPlugin struct{ *mock.Plugin }

func (HelperPlugin) Help() []help.Text {
	return []help.Text{
		{Respond: true, Command: "foo", Description: "bar"},
		{Respond: false, Command: "baz", Description: "bar"},
	}
}

type mockChat struct {
	blocking chan struct{}
	*mock.Chat
}

func newMockChat() *mockChat {
	m := &mockChat{Chat: mock.NewChat(), blocking: make(chan struct{})}
	m.Name = "bot"
	m.MessageChan = make(chan bot.Message, 1)
	m.SendFunc = func(msg bot.Message) error {
		fmt.Println(msg.Text)
		close(m.blocking)
		return nil
	}
	return m
}

func (m *mockChat) SendMessage(text string) {
	m.MessageChan <- bot.Message{Text: text}
	close(m.MessageChan)
}

func (m *mockChat) Wait() { <-m.blocking }

func Example() {
	chat := newMockChat()
	go chat.SendMessage("@bot help")

	bot.New(
		chat,
		help.New(),
		&HelperPlugin{mock.NewPlugin()},
	).Run()

	chat.Wait()
	// Output:
	// bot help - Displays all the help commands that this bot knows about.
	// bot help <query> - Displays all help commands that match <query>.
	// bot foo - bar
	// baz - bar
}

func Example_withQuery() {
	chat := newMockChat()
	go chat.SendMessage("@bot help bar")

	bot.New(
		chat,
		help.New(),
		&HelperPlugin{mock.NewPlugin()},
	).Run()

	chat.Wait()
	// Output:
	// bot foo - bar
	// baz - bar
}
