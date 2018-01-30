package bot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestPlugin struct {
	Name   string
	loaded map[string][]string
}

func (p *TestPlugin) Load(r *Robot)   { p.loaded["loaded"] = append(p.loaded["loaded"], p.Name) }
func (p *TestPlugin) Unload(r *Robot) { p.loaded["unloaded"] = append(p.loaded["unloaded"], p.Name) }

func TestPluginAdd(t *testing.T) {
	reg := newPluginRegistry()
	in := TestPlugin{Name: "Test"}
	var out TestPlugin

	ok := reg.Get(&out)
	assert.False(t, ok, "Returns false when the plugin is missing")

	reg.Add(&in)

	ok = reg.Get(&out)
	assert.True(t, ok, "Returns true when the plugin is registered")
	assert.Equal(t, in, out, "Copies the registered plugin to the pointer")
}

func TestPluginAdd_duplicate(t *testing.T) {
	reg := newPluginRegistry()
	p1 := TestPlugin{Name: "Test"}
	p2 := TestPlugin{Name: "Test"}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Adding duplicate plugins should cause a panic")
		}
	}()
	reg.Add(&p1, &p2)
}

func TestPluginLoad(t *testing.T) {
	type TestPluginTwo struct{ TestPlugin }
	out := make(map[string][]string)
	reg := newPluginRegistry()
	p1 := TestPlugin{loaded: out, Name: "one"}
	p2 := TestPluginTwo{TestPlugin{loaded: out, Name: "two"}}

	reg.Add(&p1, &p2)
	reg.Load(&Robot{})

	assert.Equal(
		t,
		[]string{"one", "two"},
		out["loaded"],
		"Should load the plugins in order",
	)

	reg.Unload(&Robot{})
	assert.Equal(
		t,
		[]string{"two", "one"},
		out["unloaded"],
		"Should unload the plugins in reverse order",
	)
}
