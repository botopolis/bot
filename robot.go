package gobot

import (
	"fmt"

	"github.com/berfarah/gobot/brain"
	"github.com/berfarah/gobot/brain/redis"
	"github.com/nlopes/slack"
)

type Plugin func(*Robot)

type Robot struct {
	rtm *slack.RTM
	api *slack.Client

	Brain *brain.Brain

	// Bot info
	ID   string
	Name string

	// Event management
	triggers map[string]trigger
}

func New(secret string) *Robot {
	redis := redis.New("localhost:6379")
	return &Robot{
		api:      slack.New(secret),
		triggers: newTriggers(),
		Brain:    brain.New(redis),
	}
}

func (r *Robot) Load(p Plugin) { p(r) }

func (r *Robot) Connect() {
	r.rtm = r.api.NewRTM()
	go r.rtm.ManageConnection()
	r.listenEvents() // blocking
}

func (r *Robot) Disconnect() error {
	return r.rtm.Disconnect()
}

func (r *Robot) onConnect(ev *slack.ConnectedEvent) {
	u := ev.Info.User
	r.ID = u.ID
	r.Name = u.Name

	r.populateCache(
		ev.Info.Users,
		ev.Info.Channels,
		ev.Info.IMs,
	)
}

func (r *Robot) Hear(l *Hook)    { r.triggers[HearTrigger].Add(l) }
func (r *Robot) Respond(l *Hook) { r.triggers["Respond"].Add(l) }
func (r *Robot) Enter(l *Hook)   { r.triggers["Enter"].Add(l) }
func (r *Robot) Leave(l *Hook)   { r.triggers["Leave"].Add(l) }
func (r *Robot) Topic(l *Hook)   { r.triggers["Topic"].Add(l) }

func (r *Robot) runListeners(ev *slack.MessageEvent) {
	if ev.User == r.ID || ev.BotID == r.ID {
		return
	}
	fmt.Printf("Message: %v\n", ev)
	for name, t := range r.triggers {
		fmt.Println("Running Trigger: " + name)
		t.Run(r, ev)
	}
}

func (r *Robot) listenEvents() {
	for msg := range r.rtm.IncomingEvents {
		// fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// fmt.Println("Hello")

		case *slack.ConnectedEvent:
			// fmt.Printf("Connected: %v\n", ev.ConnectionCount)
			r.onConnect(ev)
		case *slack.MessageEvent:
			r.runListeners(ev)
		case *slack.RTMError:
			// fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			// fmt.Printf("Invalid credentials")
			return

		default:
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
}

func (r *Robot) Send(m Message) error {
	m.Params.AsUser = true
	if m.Params.User == "" {
		m.Params.User = r.ID
	}
	_, _, err := r.api.PostMessage(m.Channel, m.Text, m.Params)
	return err
}

func (r *Robot) React(channelID, timestamp, reaction string) error {
	return r.api.AddReaction(reaction, slack.ItemRef{
		Channel:   channelID,
		Timestamp: timestamp,
	})
}

func (r *Robot) SetTopic(channelID, topic string) error {
	_, err := r.api.SetChannelTopic(channelID, topic)
	return err
}
