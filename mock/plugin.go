package mock

import "github.com/botopolis/bot"

// Plugin is a mock bot.Plugin
type Plugin struct {
	// Called by Load
	LoadFunc func(r *bot.Robot)
	// Called by Unload
	UnloadFunc func(r *bot.Robot)
}

// NewPlugin returns a NOOP Plugin
func NewPlugin() *Plugin {
	return &Plugin{
		LoadFunc:   func(r *bot.Robot) {},
		UnloadFunc: func(r *bot.Robot) {},
	}
}

// Load delegates to Plugin.LoadFunc
func (p *Plugin) Load(r *bot.Robot) { p.LoadFunc(r) }

// Unload delegates to Plugin.UnloadFunc
func (p *Plugin) Unload(r *bot.Robot) { p.UnloadFunc(r) }
