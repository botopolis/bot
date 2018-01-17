# Slack Bot

Usage:

```go
package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/berfarah/go/bot"
)

func main() {
	var secret string
	if contents, err := ioutil.ReadFile(".secret"); err != nil {
		panic(err)
	} else {
		secret = strings.TrimSpace(string(contents))
	}

	r := bot.New(secret)
	r.Hear(&bot.Hook{
		Name:    "yoyo",
		Matcher: bot.MatchText("heyy"),
		Func: func(r *bot.Responder) error {
			fmt.Println("running thing")
			r.Send(bot.Message{Text: "ayy"})
			return nil
		},
	})
	r.Connect()
}
```
