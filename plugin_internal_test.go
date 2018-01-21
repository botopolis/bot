package gobot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestPlugin struct {
	Name  string
	State string
}

func (p *TestPlugin) Load(r *Robot)   { p.State = "Loaded" }
func (p *TestPlugin) Unload(r *Robot) { p.State = "Unloaded" }

func TestPluginAdd(t *testing.T) {
	reg := newPluginRegistry()
	in := TestPlugin{Name: "Test", State: "Loading"}
	var out TestPlugin

	ok := reg.Get(&out)
	assert.False(t, ok, "Returns false when the plugin is missing")

	reg.Add(&in)

	ok = reg.Get(&out)
	assert.True(t, ok, "Returns true when the plugin is registered")
	assert.Equal(t, in, out, "Copies the registered plugin to the pointer")
}

func TestPluginLoad(t *testing.T) {
	reg := newPluginRegistry()
	in := TestPlugin{Name: "Test", State: "Loading"}
	var out TestPlugin

	reg.Add(&in)

	reg.Load(&Robot{})
	reg.Get(&out)
	assert.Equal(t, "Loaded", out.State)

	reg.Unload(&Robot{})
	reg.Get(&out)
	assert.Equal(t, "Unloaded", out.State)
}
