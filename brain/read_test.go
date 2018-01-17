package brain_test

import (
	"testing"

	"github.com/berfarah/gobot/brain"
	"github.com/berfarah/gobot/brain/memory"
	"github.com/stretchr/testify/assert"
)

type Readable struct {
	id   string
	Name string
}

func (r Readable) ID() string        { return r.id }
func (r Readable) Namespace() string { return "normal" }

func TestRead(t *testing.T) {
	assert := assert.New(t)
	r := Readable{id: "5", Name: "foo"}
	m := memory.New()
	m.Cache["normal:5"] = r
	b := brain.New(m)

	v := Readable{id: "5"}
	b.Read(&v)

	assert.Equal(r, v)
}

type ReadableIndexable struct {
	id   string
	Name string
}

func (r ReadableIndexable) ID() string        { return r.id }
func (r ReadableIndexable) Namespace() string { return "normal" }
func (r ReadableIndexable) Indices() map[string]string {
	return map[string]string{"name": r.Name}
}

func TestReadWithIndex(t *testing.T) {
	assert := assert.New(t)
	r := ReadableIndexable{id: "5", Name: "foo"}
	m := memory.New()
	m.Cache["normal:5"] = r
	m.Cache["normal_by_name:foo"] = "5"
	b := brain.New(m)

	v := ReadableIndexable{Name: "foo"}
	b.Read(&v)

	assert.Equal(r, v)
}
