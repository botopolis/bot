package gobot

import (
	"fmt"

	"github.com/nlopes/slack"
)

func (r *Robot) listenEvents() {
	r.rtm = r.api.NewRTM()
	go r.rtm.ManageConnection()
	for msg := range r.rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
		case *slack.ConnectedEvent:
			r.onConnect(ev)
			r.Logger.Debugf("Connected as %s: %d", r.ID, ev.ConnectionCount)
		case *slack.MessageEvent:
			r.runListeners(ev)
		case *slack.RTMError:
			r.Logger.Errorf("RTM Error: %s", ev.Error())
		case *slack.InvalidAuthEvent:
			r.Logger.Error("Slack: Invalid Credentials")
			return
		default:
			r.Logger.Debugf("Unhandled Slack event: %s", msg.Type)
		}
	}
}

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
