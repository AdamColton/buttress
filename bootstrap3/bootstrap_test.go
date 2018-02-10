package bootstrap3

import (
	"fmt"
	"github.com/adamcolton/buttress/bootstrap3/csscontext"
	"github.com/adamcolton/buttress/bootstrap3/csssize"
	"github.com/adamcolton/buttress/html"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNav(t *testing.T) {
	d := NewDocument("test")
	nav := NewNav("Gothic Testing", "#home")
	nav.Add(Right, "register", "Register", "").SetHref("/user/create")
	nav.Add(Left, "home", "Home", "home").SetHref("/")

	d.StandardNav(nav)
	s := d.String()
	for _, link := range CSSLinks {
		assert.Contains(t, s, link)
	}
	for _, link := range ScriptLinks {
		assert.Contains(t, s, link)
	}
}

func TestPanel(t *testing.T) {
	p := NewPanel("Testing", html.NewText("this is a test"))
	assert.Contains(t, html.String(p.Render()), "Testing")
}

func TestColClass(t *testing.T) {
	c := NewColClass()
	c.SetSize(5, csssize.Small())

	size, ok := c.Size(csssize.Small())
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	size, ok = c.Size(csssize.Large())
	assert.False(t, ok)
	assert.Equal(t, 5, size)

	assert.Equal(t, "col-sm-5", c.String())
	c.SetSize(5, csssize.ExtraSmall())
	assert.Equal(t, "col-xs-5", c.String())
	c.SetSize(8, csssize.ExtraSmall())
	assert.Equal(t, "col-xs-8 col-sm-5", c.String())
	c.SetHidden(true, csssize.Medium())
	assert.Equal(t, "col-xs-8 col-sm-5 hidden-md", c.String())
	c.SetHidden(true, csssize.Small())
	c.SetHidden(true, csssize.ExtraSmall())
	assert.Equal(t, "col-xs-8 col-sm-5 visible-lg", c.String())
}

func TestRow(t *testing.T) {
	r := NewRow(nil)
	c := NewColClass()
	c.SetSize(4, csssize.ExtraSmall())
	c.SetSize(3, csssize.Medium())

	for i := 0; i < 13; i++ {
		node := html.NewTag("div")
		node.AddChildren(html.NewText(fmt.Sprintf("Cell %d", i)))
		r.AddNodes(c, node)
	}

	out := r.Render()
	s := html.String(out)
	assert.Contains(t, s, `<div class="clearfix visible-xs-block visible-sm-block visible-md-block visible-lg-block"></div>`)
}

func TestContainer(t *testing.T) {
	cc := NewColClass()
	cc.SetSize(12, csssize.ExtraSmall())
	cc.SetSize(6, csssize.Small())
	cc.SetSize(3, csssize.Medium())
	cc.SetSize(1, csssize.Large())

	c := NewFluidContainer()

	for y := 0; y < 5; y++ {
		for x := 0; x < 4; x++ {
			node := html.NewTag("div")
			node.AddChildren(html.NewText(fmt.Sprintf("Cell %d, %d", x, y)))
			c.AddNodes(cc, node)
		}
	}

	out := c.Render()
	s := html.String(out)
	assert.Contains(t, s, `<div class="container-fluid">`)
	// Todo: for clear fix, use hidden/visible correctly
	assert.Contains(t, s, `<div class="clearfix visible-xs-block visible-sm-block visible-md-block visible-lg-block"></div>`)
}

func TestDemo(t *testing.T) {
	d := NewDocument("Bootstrap Demo")

	nav := NewNav("Gothic Testing", "#home")
	nav.Add(Right, "register", "Register", "").SetHref("#register")
	nav.Add(Left, "home", "Home", "home").SetHref("#home")
	d.StandardNav(nav)

	f := NewForm(nil)
	f.Style.Label.SetSize(4, csssize.Small())
	f.Style.Input.SetSize(5, csssize.Small())
	f.Style.Input.SetSize(4, csssize.Medium())
	f.Style.Input.SetSize(3, csssize.Large())
	f.AddInputTag("text", "Name", "name")
	f.AddInputTag("text", "Age", "age")
	f.AddInputTag("text", "Role", "role")
	f.AddButtons("Save", "save", "floppy-o", csscontext.Success()).AddButton("Cancel", "cancel", "", csscontext.Danger())

	panel := NewPanel("Hello", f.Render())
	panel.ContextClass = csscontext.Primary()

	cc := NewColClass()
	cc.SetSize(12, csssize.ExtraSmall())
	cc.SetSizeOffset(10, 1, csssize.Medium())

	c := NewContainer()
	c.AddNodes(cc, panel.Render())

	d.AddChildren(c.Render())

	s := d.String()
	assert.Contains(t, s, "<html>")
}

func TestForm(t *testing.T) {
	fs := NewFormStyle()
	fs.Label.SetSizeOffset(3, 2, csssize.ExtraSmall())
	fs.Input.SetSize(4, csssize.ExtraSmall())

	f := &Form{
		Style: fs,
	}
	f.AddInputTag("text", "Name", "name")
	f.AddInputTag("text", "Age", "age")
	f.AddInputTag("text", "Role", "role")
	f.AddText("TextTest", "here is some text")
	f.AddConfirmPassword("password")
	f.AddSelect("SelectTest", "selectMe", []SelectOption{
		{"foo", "FOO"},
		{"bar", "Bar"},
	})
	f.AddButton("Save", "user")

	s := html.String(f.Render())
	assert.Contains(t, s, `<form class="form-horizontal">`)

	println(s)
}
