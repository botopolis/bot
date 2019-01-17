package help

import (
	"fmt"
	"strings"
	"sync"

	"github.com/botopolis/bot"
)

// Plugin is the help Plugin for botopolis/bot. It implements
// bot.Plugin and help.Helper.
type Plugin struct {
	bot   *bot.Robot
	once  sync.Once
	cache []string
}

// New creates a new instance of the help plugin.
func New() bot.Plugin { return &Plugin{} }

// Provider is the interface a Plugin is expected to implement
// in order to provide help text.
type Provider interface {
	Help() []Text
}

// Text is a representation of what is returned to the user.
// Example:
// Text{Respond: true, Command: "foo", Description: "bar"}
// turns into
// "<robot.Username> foo - bar"
type Text struct {
	// Respond is true if the command needs to be sent at a bot.
	// Example: @bot <command>
	Respond bool
	// The command for which we are writing help text.
	// Example: create link <name> <link>
	Command string
	// The description of the command.
	// Example: Creates a link with the <name> short URL to the <link>.
	Description string
}

func (p Text) format(r *bot.Robot) string {
	var at string
	if p.Respond {
		at = fmt.Sprintf("%s ", r.Username())
	}

	return fmt.Sprintf("%s%s - %s", at, p.Command, p.Description)
}

// Load starts the help.Plugin, starts listening for help commands.
func (p *Plugin) Load(r *bot.Robot) {
	p.bot = r

	r.Respond(bot.Regexp(`help\s*(.*)`), func(r bot.Responder) error {
		p.loadCache()

		return r.Send(bot.Message{
			Text: p.buildResponse(r.Match[1]),
		})
	})
}

func (p *Plugin) loadCache() {
	p.once.Do(func() {
		for _, plugin := range p.bot.ListPlugins() {
			if helper, ok := plugin.(Provider); ok {
				for _, helpText := range helper.Help() {
					p.cache = append(p.cache, helpText.format(p.bot))
				}
			}
		}
	})
}

func (p *Plugin) buildResponse(query string) string {
	if query == "" {
		return strings.Join(p.cache, "\n")
	}

	var partialMatch []string
	for _, text := range p.cache {
		if strings.Contains(strings.ToLower(text), strings.ToLower(query)) {
			partialMatch = append(partialMatch, text)
		}
	}
	return strings.Join(partialMatch, "\n")
}

// Help provides helper text and implements the help.Helper interface.
func (p *Plugin) Help() []Text {
	return []Text{
		{
			Respond:     true,
			Command:     "help",
			Description: "Displays all the help commands that this bot knows about.",
		},
		{
			Respond:     true,
			Command:     "help <query>",
			Description: "Displays all help commands that match <query>.",
		},
	}
}

var _ Provider = &Plugin{}
