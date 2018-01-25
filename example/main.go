package main

import (
	"os"

	"github.com/berfarah/gobot"
	"github.com/berfarah/gobot-slack"
	"github.com/berfarah/gobot-store-redis"
)

func main() {
	gobot.New(
		slack.New(os.Getenv("SLACK_TOKEN")),
		redis.New(os.Getenv("REDIS_URL")),
		welcomer{},
		webhook{},
	).Run()
}
