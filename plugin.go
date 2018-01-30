package bot

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
	registry  map[reflect.Type]Plugin
	loadOrder []reflect.Type
}

func newPluginRegistry() *pluginRegistry {
	return &pluginRegistry{
		registry:  make(map[reflect.Type]Plugin),
		loadOrder: make([]reflect.Type, 0),
	}
}

func (reg *pluginRegistry) Add(plugins ...Plugin) {
	var (
		duplicate []string
	)
	for _, p := range plugins {
		t := reflect.TypeOf(p)
		if _, ok := reg.registry[t]; ok {
			duplicate = append(duplicate, t.PkgPath()+":"+t.Name())
			continue
		}
		reg.registry[t] = p
		reg.loadOrder = append(reg.loadOrder, t)
	}

	if len(duplicate) > 0 {
		panic(fmt.Sprintf("Warning - Duplicate plugins added: %v", duplicate))
	}
}

func (reg *pluginRegistry) Get(p Plugin) bool {
	t := reflect.TypeOf(p)
	if plugin, ok := reg.registry[t]; ok {
		copyInterface(plugin, p)
		return true
	}
	return false
}

func (reg *pluginRegistry) Load(r *Robot) {
	for _, t := range reg.loadOrder {
		if p, ok := reg.registry[t]; ok {
			p.Load(r)
		}
	}
}

func (reg *pluginRegistry) Unload(r *Robot) {
	for i := (len(reg.loadOrder) - 1); i >= 0; i-- {
		t := reg.loadOrder[i]
		if p, ok := reg.registry[t]; ok {
			if u, ok := p.(unloader); ok {
				u.Unload(r)
			}
		}
	}
}
