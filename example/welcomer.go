package main

import "github.com/botopolis/bot"

type welcomer struct{}

func (w welcomer) Load(r *bot.Robot) {
	r.Hear(bot.Regexp("welcome (\\w*)"), func(r bot.Responder) error {
		return r.Send(bot.Message{
			Text: "All hail " + r.Match[1] + ", our new overlord",
		})
	})
}
