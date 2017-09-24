package bootstrap3

import (
	"github.com/adamcolton/buttress/bootstrap3/csscontext"
	"github.com/adamcolton/buttress/html"
	"github.com/adamcolton/buttress/html/builder"
)

// https://getbootstrap.com/docs/3.3/components/#panels

func NewPanel(title string, contents html.Node) *Panel {
	return &Panel{
		ContextClass: csscontext.Default(),
		Title:        title,
		Contents:     contents,
	}
}

type Panel struct {
	ContextClass csscontext.CSSContext
	Title        string
	Footer       string
	Contents     html.Node
}

func (p *Panel) Render() html.Node {
	b := builder.New().Tag("div", "class", "panel panel-"+p.ContextClass.String())
	ret := b.Cur()
	if p.Title != "" {
		b.Tag("div", "class", "panel-heading").Text(p.Title).Close()
	}
	b.Tag("div", "class", "panel-body").
		AddChildren(p.Contents)
	return ret
}
