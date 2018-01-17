package gobot

import (
	"fmt"
	"regexp"
)

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
		return fmt.Errorf("Hook Error in %s: %s", h.Name, err.Error())
	}
	return nil
}

func MatchUser(u string) Matcher {
	return func(r *Responder) bool {
		return r.Message.User == u
	}
}

func MatchChannel(c string) Matcher {
	return func(r *Responder) bool {
		return r.Message.Channel == c
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
