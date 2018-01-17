# Slack Bot

The vision for this project is:

Create a slack interface for hooking in multiple callbacks, an abstraction to
make interacting with it easy

Create an interface for configuration data:
- Properties for slack users (JSON)
- Properties for channels (JSON)
- Ability to have indexed properties

Add other integrations like GitHub and CI

Add building blocks for "typical interactions", whatever those turn out to be

## Usage

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
