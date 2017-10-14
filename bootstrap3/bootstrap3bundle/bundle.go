package bootstrap3bundle

import (
	"github.com/adamcolton/buttress/bootstrap3"
	"github.com/adamcolton/buttress/bootstrap3/bootstrap3model"
	"github.com/adamcolton/buttress/bootstrap3/csscontext"
	"github.com/adamcolton/buttress/bootstrap3/csssize"
	"github.com/adamcolton/buttress/html"
)

type Bundle struct {
	SinglePanelClass struct {
		Col     bootstrap3.ColClass
		Context csscontext.CSSContext
	}
	FormStyle       *bootstrap3.FormStyle
	FormBuilder     *bootstrap3model.FormBuilder
	Nav             *bootstrap3.Nav
	Title           string
	CSS             []string
	Scripts         []string
	DocumentBuilder func(d *bootstrap3.Document)
}

var ConfirmPassword = bootstrap3model.TypeMap{"password": bootstrap3model.ConfirmPassword}

func New(title string) *Bundle {
	b := &Bundle{
		FormStyle: bootstrap3.NewFormStyle(),
		Nav:       bootstrap3.NewNav(title, ""),
		Title:     title,
	}

	b.SinglePanelClass.Col = bootstrap3.NewColClass()
	b.SinglePanelClass.Col.SetSize(12, csssize.Medium())
	b.SinglePanelClass.Col.SetSizeOffset(10, 1, csssize.Large())
	b.SinglePanelClass.Context = csscontext.Default()

	b.FormStyle.Label.SetSize(5, csssize.Medium())
	b.FormStyle.Input.SetSize(4, csssize.Medium())
	b.FormBuilder = bootstrap3model.NewFormBuilder(b.FormStyle)

	b.DocumentBuilder = b.DefaultBuilder
	return b
}

func (b *Bundle) AddScripts(scripts ...string) {
	b.Scripts = append(b.Scripts, scripts...)
}

func (b *Bundle) AddCSS(css ...string) {
	b.CSS = append(b.CSS, css...)
}

func (b *Bundle) Document() *bootstrap3.Document {
	d := bootstrap3.NewDocument(b.Title)
	d.CSSLinks(b.CSS...)
	d.ScriptLinks(b.Scripts...)
	if b.DocumentBuilder != nil {
		b.DocumentBuilder(d)
	}
	return d
}

func (b *Bundle) DefaultBuilder(d *bootstrap3.Document) {
	d.StandardNav(b.Nav)
}

func (b *Bundle) SinglePanel(title string, contents html.Node) *bootstrap3.Container {
	panel := bootstrap3.NewPanel(title, contents)
	panel.ContextClass = b.SinglePanelClass.Context

	c := bootstrap3.NewContainer()
	c.AddNodes(b.SinglePanelClass.Col, panel.Render())
	return c
}

func (b *Bundle) SinglePanelDocument(title string, contents html.Node) *bootstrap3.Document {
	d := b.Document()
	d.AddChildren(b.SinglePanel(title, contents).Render())
	return d
}

func (b *Bundle) Form() *bootstrap3.Form {
	return bootstrap3.NewForm(b.FormStyle)
}
