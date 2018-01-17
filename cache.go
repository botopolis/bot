package gobot

import (
	"errors"

	"github.com/nlopes/slack"
)

type User struct{ slack.User }

func (u User) ID() string        { return u.User.ID }
func (u User) Namespace() string { return "slack:user" }
func (u User) Indices() map[string]string {
	return map[string]string{"username": u.Name}
}

type Channel struct{ slack.Channel }

func (c Channel) ID() string        { return c.Channel.ID }
func (c Channel) Namespace() string { return "slack:channel" }
func (c Channel) Indices() map[string]string {
	return map[string]string{"name": c.Name}
}

type DM struct {
	Username     string
	UserID       string
	Conversation string
}

func (d DM) ID() string        { return d.Conversation }
func (d DM) Namespace() string { return "slack:dm" }
func (d DM) Indices() map[string]string {
	return map[string]string{"username": d.Username, "user_id": d.UserID}
}

func (r *Robot) updateCache() error {
	users, err := r.api.GetUsers()
	if err != nil {
		return errors.New("Couldn't fetch users")
	}

	channels, err := r.api.GetChannels(true)
	if err != nil {
		return errors.New("Couldn't fetch channels")
	}

	dms, err := r.api.GetIMChannels()
	if err != nil {
		return errors.New("Couldn't fetch DMs")
	}

	return r.populateCache(users, channels, dms)
}

func (r *Robot) populateCache(users []slack.User, channels []slack.Channel, dms []slack.IM) error {
	userIDToName := make(map[string]string)
	for _, u := range users {
		r.Brain.Write(User{u})
		userIDToName[u.ID] = u.Name
	}

	for _, ch := range channels {
		r.Brain.Write(Channel{ch})
	}

	for _, dm := range dms {
		if name, ok := userIDToName[dm.User]; ok {
			r.Brain.Write(DM{
				Conversation: dm.ID,
				UserID:       dm.User,
				Username:     name,
			})
		}
	}

	return nil
}
