package gobot

import (
	"fmt"
	"reflect"
)

// Plugin is anything that can be installed into Robot
type Plugin interface {
	Load(*Robot)
}

type unloader interface {
	Unload(*Robot)
}

type pluginRegistry struct {
	registry map[reflect.Type]Plugin
}

func newPluginRegistry() *pluginRegistry {
	return &pluginRegistry{registry: make(map[reflect.Type]Plugin)}
}

func (reg *pluginRegistry) Add(plugins ...Plugin) {
	var duplicate []string
	for _, p := range plugins {
		t := reflect.TypeOf(p)
		if _, ok := reg.registry[t]; ok {
			duplicate = append(duplicate, t.PkgPath()+":"+t.Name())
		}
		reg.registry[t] = p
	}

	if len(duplicate) > 0 {
		panic(fmt.Sprintf("Warning - Duplicate plugins added: %v", duplicate))
	}
}

func (reg *pluginRegistry) Get(p Plugin) bool {
	t := reflect.TypeOf(p)
	if plugin, ok := reg.registry[t]; ok {
		if err := copyInterface(plugin, p); err != nil {
			return false
		}
		return true
	}
	return false
}

func (reg *pluginRegistry) Load(r *Robot) int {
	for _, p := range reg.registry {
		p.Load(r)
	}
	return len(reg.registry)
}

func (reg *pluginRegistry) Unload(r *Robot) int {
	var count int
	for _, p := range reg.registry {
		if u, ok := p.(unloader); ok {
			u.Unload(r)
			count++
		}
	}
	return count
}
