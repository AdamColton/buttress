package mutate

import (
	"github.com/adamcolton/buttress/html"
	"github.com/adamcolton/buttress/html/query"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMutate(t *testing.T) {
	before := html.NewTag("div", "class", "top")
	before.AddChildren(html.NewText("This is a test"))

	q, err := query.Selector(".top")
	assert.NoError(t, err)
	var m MutateFunc = func(root html.Node) html.Node {
		q.Query(root).AppendClass("mutated")
		return root
	}

	mc := MutateChain{m}
	after := mc.Mutate(before)

	assert.Equal(t, `<div class="top mutated">This is a test</div>`, html.String(after))
}

func TestAppendTags(t *testing.T) {
	before := html.NewTag("div", "class", "top")
	before.AddChildren(html.NewText("This is a test"))

	m, err := AppendTags(".top", "class", "mutated", "foo", "bar")
	assert.NoError(t, err)
	after := m.Mutate(before)

	assert.Equal(t, `<div class="top mutated" foo="bar">This is a test</div>`, html.String(after))
}

func TestAppendClass(t *testing.T) {
	before := html.NewTag("div", "class", "top")
	before.AddChildren(html.NewText("This is a test"))

	m, err := AppendClass(".top", "mutated")
	assert.NoError(t, err)
	after := m.Mutate(before)

	assert.Equal(t, `<div class="top mutated">This is a test</div>`, html.String(after))
}
