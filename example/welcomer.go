package main

import "github.com/berfarah/gobot"

type welcomer struct{}

func (w welcomer) Load(r *gobot.Robot) {
	r.Hear(gobot.Regexp("welcome (\\w*)"), func(r gobot.Responder) error {
		return r.Send(gobot.Message{
			Text: "All hail " + r.Match[1] + ", our new overlord",
		})
	})
}
