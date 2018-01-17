package gobot

import (
	"regexp"

	"github.com/nlopes/slack"
)

const (
	HearTrigger    = "Hear"
	RespondTrigger = "Respond"
	EnterTrigger   = "Enter"
	LeaveTrigger   = "Leave"
	TopicTrigger   = "Topic"
)

type trigger interface {
	Add(*Hook)
	Run(*Robot, *slack.MessageEvent) bool
}

type messageTrigger struct {
	hooks   []*Hook
	subType string
}
type respondTrigger struct{ messageTrigger }

func newMessageTrigger(s string) *messageTrigger {
	return &messageTrigger{subType: s}
}

func newRespondTrigger(s string) *respondTrigger {
	return &respondTrigger{messageTrigger{subType: s}}
}

func (mt *messageTrigger) Add(h *Hook) {
	if len(mt.hooks) == 0 {
		mt.hooks = []*Hook{h}
		return
	}
	mt.hooks = append(mt.hooks, h)
}

func (mt *messageTrigger) Run(r *Robot, msg *slack.MessageEvent) bool {
	if msg.SubType != mt.subType {
		return false
	}
	for _, h := range mt.hooks {
		h.Run(newResponder(r, msg))
	}
	return true
}

func (rt *respondTrigger) Run(r *Robot, msg *slack.MessageEvent) bool {
	if msg.SubType != rt.subType {
		return false
	}

	reg, err := regexp.Compile("<@" + r.ID + ">")
	if err != nil || !reg.MatchString(msg.Text) {
		return false
	}
	for _, h := range rt.hooks {
		h.Run(newResponder(r, msg))
	}
	return true
}

func newTriggers() map[string]trigger {
	return map[string]trigger{
		HearTrigger:    newMessageTrigger(""),
		RespondTrigger: newRespondTrigger(""),
		EnterTrigger:   newMessageTrigger("channel_join"),
		LeaveTrigger:   newMessageTrigger("channel_leave"),
		TopicTrigger:   newMessageTrigger("channel_topic"),
	}
}
