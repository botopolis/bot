package gobot

import (
	"regexp"
)

const (
	HearTrigger    = "Hear"
	RespondTrigger = "Respond"
	EnterTrigger   = "Enter"
	LeaveTrigger   = "Leave"
	TopicTrigger   = "Topic"
)

func newTriggers() map[string][]*Hook {
	return make(map[string][]*Hook)
}

func (r *Robot) trigger(on string, h *Hook) {
	r.triggers[on] = append(r.triggers[on], h)
}

func (r *Robot) Hear(l *Hook)    { r.trigger(HearTrigger, l) }
func (r *Robot) Respond(l *Hook) { r.trigger(RespondTrigger, l) }
func (r *Robot) Enter(l *Hook)   { r.trigger(EnterTrigger, l) }
func (r *Robot) Leave(l *Hook)   { r.trigger(LeaveTrigger, l) }
func (r *Robot) Topic(l *Hook)   { r.trigger(TopicTrigger, l) }

type Hook struct {
	Name    string
	Func    func(r *Responder) error
	Matcher func(r *Responder) bool
}
type Matcher func(r *Responder) bool

func (h *Hook) Run(r *Responder) error {
	if h.Matcher != nil && !h.Matcher(r) {
		return nil
	}
	if err := h.Func(r); err != nil {
		return err
	}
	return nil
}

func MatchUser(u string) Matcher {
	return func(r *Responder) bool {
		return r.Message.User == u
	}
}

func MatchRoom(c string) Matcher {
	return func(r *Responder) bool {
		return r.Message.Room == c
	}
}

func MatchText(t string) Matcher {
	return func(r *Responder) bool {
		return r.Message.Text == t
	}
}

func MatchRegexp(expression string) Matcher {
	return func(r *Responder) bool {
		reg, err := regexp.Compile(expression)
		if err == nil && reg.MatchString(r.Message.Text) {
			r.Match = reg.FindStringSubmatch(r.Message.Text)
			return true
		}
		return false
	}
}
