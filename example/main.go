package main

import (
	"os"

	"github.com/botopolis/bot"
	"github.com/botopolis/redis"
	"github.com/botopolis/slack"
)

func main() {
	bot.New(
		slack.New(os.Getenv("SLACK_TOKEN")),
		redis.New(os.Getenv("REDIS_URL")),
		welcomer{},
		webhook{},
	).Run()
}
