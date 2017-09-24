package query

import (
	"github.com/adamcolton/buttress/html"
	"github.com/adamcolton/buttress/html/parsehtml"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMatches(t *testing.T) {
	tag := html.NewTag("div", "class", "foo glorp", "id", "testing")
	assert.True(t, checkTag("div").check(tag))
	assert.False(t, checkTag("ul").check(tag))
	assert.True(t, checkClass("foo").check(tag))
	assert.True(t, checkClass("glorp").check(tag))
	assert.False(t, checkClass("bar").check(tag))
	assert.True(t, checkID("testing").check(tag))
	assert.False(t, checkID("test").check(tag))
}

func TestSelect(t *testing.T) {
	i := parsehtml.Must("<i>bar</i>")
	p1 := parsehtml.Must("<p>paragraph 1</p>")
	p2 := parsehtml.Must(`<p class="foo">paragraph 2</p>`)
	p3 := parsehtml.Must(`<p class="foo">paragraph 3</p>`)
	div := html.NewTag("div")
	div.AddChildren(i, p1, html.NewText("interrupt"), html.NewFragment(p2, p3))
	p4 := html.NewTag("p", "class", "foo")
	p4.AddChildren(html.NewText("paragraph 4"))
	body := html.NewTag("body")
	body.AddChildren(div, html.NewVoidTag("hr"), p4)

	title := parsehtml.Must(`<title>This is a test</title>`)
	head := html.NewTag("head")
	head.AddChildren(title)

	html := html.NewTag("html")
	html.AddChildren(head, body)
	root := html.NewFragment(html.NewDoctype("html"), html)

	s := selectors{
		&selector{
			checkers: []nodeChecker{checkTag("div")},
		},
	}
	matches := s.QueryAll(root)
	if assert.Len(t, matches, 1) {
		assert.Equal(t, div, matches[0])
	}

	s = selectors{
		&selector{
			checkers: []nodeChecker{checkTag("p"), checkClass("foo")},
		},
	}
	matches = s.QueryAll(root)
	if assert.Len(t, matches, 3) {
		assert.Equal(t, p2, matches[0])
		assert.Equal(t, p3, matches[1])
		assert.Equal(t, p4, matches[2])
	}

	s = selectors{
		&selector{
			checkers: []nodeChecker{checkTag("div")},
			next: &selector{
				checkers: []nodeChecker{checkTag("p"), checkClass("foo")},
			},
			nextLoc: newDescendant,
		},
	}

	sel, err := Selector("div p.foo, title")
	assert.NoError(t, err)
	assert.NotNil(t, sel)
	matches = sel.QueryAll(root)
	if assert.Len(t, matches, 3) {
		assert.Equal(t, title, matches[0])
		assert.Equal(t, p2, matches[1])
		assert.Equal(t, p3, matches[2])
	}

	sel, err = Selector("div>p")
	assert.NoError(t, err)
	assert.NotNil(t, sel)
	matches = sel.QueryAll(root)
	if assert.Len(t, matches, 3) {
		assert.Equal(t, p1, matches[0])
		assert.Equal(t, p2, matches[1])
		assert.Equal(t, p3, matches[2])
	}

	sel, err = Selector("i+p")
	assert.NoError(t, err)
	assert.NotNil(t, sel)
	matches = sel.QueryAll(root)
	if assert.Len(t, matches, 1) {
		assert.Equal(t, p1, matches[0])
	}

	sel, err = Selector("i~p")
	assert.NoError(t, err)
	assert.NotNil(t, sel)
	matches = sel.QueryAll(root)
	if assert.Len(t, matches, 3) {
		assert.Equal(t, p1, matches[0])
		assert.Equal(t, p2, matches[1])
		assert.Equal(t, p3, matches[2])
	}
	match := sel.Query(root)
	assert.Equal(t, p1, match)
	paths := sel.QueryAllPaths(root)
	if assert.Len(t, paths, 3) {
		matches := paths.Query(root)
		assert.Equal(t, p1, matches[0])
		assert.Equal(t, p2, matches[1])
		assert.Equal(t, p3, matches[2])
	}
	path := sel.QueryPath(root)
	assert.Equal(t, p1, path.Query(root))
}
