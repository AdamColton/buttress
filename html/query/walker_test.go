package query

import (
	"github.com/adamcolton/buttress/html"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWalker(t *testing.T) {
	p1 := html.NewTag("p")
	p1.AddChildren(html.NewText("paragraph 1"))
	p2 := html.NewTag("p")
	p2.AddChildren(html.NewText("paragraph 2"))
	p3 := html.NewTag("p")
	p3.AddChildren(html.NewText("paragraph 3"))
	div := html.NewTag("div")
	div.AddChildren(p1, html.NewFragment(p2, p3))
	body := html.NewTag("body")
	body.AddChildren(div, html.NewVoidTag("hr"))

	title := html.NewTag("title")
	title.AddChildren(html.NewText("This is a test"))
	head := html.NewTag("head")
	head.AddChildren(title)

	html := html.NewTag("html")
	html.AddChildren(head, body)
	root := html.NewFragment(html.NewDoctype("html"), html)

	assert.NotNil(t, root)

	expected := []struct {
		node html.Node
		loc  *Location
	}{
		{
			node: root,
			loc:  &Location{},
		},
		{
			node: root.Child(0),
			loc: &Location{
				Path: NewPath(0),
				Tag:  []int{0},
				Node: []int{0},
			},
		},
		{
			node: html,
			loc: &Location{
				Path: NewPath(1),
				Tag:  []int{0},
				Node: []int{1},
			},
		},
		{
			node: head,
			loc: &Location{
				Path: NewPath(1, 0),
				Tag:  []int{0, 0},
				Node: []int{1, 0},
			},
		},
		{
			node: title,
			loc: &Location{
				Path: NewPath(1, 0, 0),
				Tag:  []int{0, 0, 0},
				Node: []int{1, 0, 0},
			},
		},
		{
			node: title.Child(0),
			loc: &Location{
				Path: NewPath(1, 0, 0, 0),
				Tag:  []int{0, 0, 0, 0},
				Node: []int{1, 0, 0, 0},
			},
		},
		{
			node: body,
			loc: &Location{
				Path: NewPath(1, 1),
				Tag:  []int{0, 1},
				Node: []int{1, 1},
			},
		},
		{
			node: div,
			loc: &Location{
				Path: NewPath(1, 1, 0),
				Tag:  []int{0, 1, 0},
				Node: []int{1, 1, 0},
			},
		},
		{
			node: p1,
			loc: &Location{
				Path: NewPath(1, 1, 0, 0),
				Tag:  []int{0, 1, 0, 0},
				Node: []int{1, 1, 0, 0},
			},
		},
		{
			node: p1.Child(0),
			loc: &Location{
				Path: NewPath(1, 1, 0, 0, 0),
				Tag:  []int{0, 1, 0, 0, 0},
				Node: []int{1, 1, 0, 0, 0},
			},
		},
		{
			node: div.Child(1),
			loc: &Location{
				Path: NewPath(1, 1, 0, 1),
				Tag:  []int{0, 1, 0, 1},
				Node: []int{1, 1, 0, 1},
			},
		},
		{
			node: p2,
			loc: &Location{
				Path: NewPath(1, 1, 0, 1, 0),
				Tag:  []int{0, 1, 0, 1},
				Node: []int{1, 1, 0, 1},
			},
		},
		{
			node: p2.Child(0),
			loc: &Location{
				Path: NewPath(1, 1, 0, 1, 0, 0),
				Tag:  []int{0, 1, 0, 1, 0},
				Node: []int{1, 1, 0, 1, 0},
			},
		},
		{
			node: p3,
			loc: &Location{
				Path: NewPath(1, 1, 0, 1, 1),
				Tag:  []int{0, 1, 0, 2},
				Node: []int{1, 1, 0, 2},
			},
		},
		{
			node: p3.Child(0),
			loc: &Location{
				Path: NewPath(1, 1, 0, 1, 1, 0),
				Tag:  []int{0, 1, 0, 2, 0},
				Node: []int{1, 1, 0, 2, 0},
			},
		},
		{
			node: body.Child(1),
			loc: &Location{
				Path: NewPath(1, 1, 1),
				Tag:  []int{0, 1, 1},
				Node: []int{1, 1, 1},
			},
		},
	}
	var i int
	visiter := func(node html.Node, location *Location) error {
		if i >= len(expected) {
			t.Error("Too many nodes")
			return nil
		}
		e := expected[i]
		i++
		assert.Equal(t, e.node, node)
		assert.EqualValues(t, e.loc, location)
		return nil
	}

	Walk(root, visiter)
}
