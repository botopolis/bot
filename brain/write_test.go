package brain_test

import (
	"testing"

	"github.com/berfarah/gobot/brain"
	"github.com/berfarah/gobot/brain/memory"
	"github.com/stretchr/testify/assert"
)

type Writeable struct {
	id   string
	Name string
}

func (w Writeable) ID() string        { return w.id }
func (w Writeable) Namespace() string { return "normal" }

func TestWrite(t *testing.T) {
	assert := assert.New(t)
	w := Writeable{id: "5", Name: "foo"}
	m := memory.New()
	b := brain.New(m)
	b.Write(w)

	v, ok := m.Cache["normal:5"]

	assert.True(ok, "should add the interface under the namespace")
	assert.Equal(w, v)
}

type WriteableIndexable struct {
	id   string
	Name string
}

func (w WriteableIndexable) ID() string        { return w.id }
func (w WriteableIndexable) Namespace() string { return "normal" }
func (w WriteableIndexable) Indices() map[string]string {
	return map[string]string{"name": w.Name}
}

func TestWriteWithIndex(t *testing.T) {
	assert := assert.New(t)
	w := WriteableIndexable{id: "5", Name: "foo"}
	m := memory.New()
	b := brain.New(m)
	b.Write(w)

	v, ok := m.Cache["normal_by_name:foo"]

	assert.True(ok, "should add the interface under the namespace")
	assert.Equal("5", v)
}
