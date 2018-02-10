package builder

import (
	"github.com/adamcolton/buttress/html"
	"io"
)

type Document struct {
	header  *html.Tag
	body    *html.Fragment
	Scripts *html.Fragment
	doc     *html.Fragment
}

func NewDocument(title string, attrs ...string) *Document {
	doc := &Document{
		doc: html.NewFragment(html.NewDoctype("html")),
	}
	doc.header = html.NewTag("header")
	if title != "" {
		titleTag := html.NewTag("title")
		titleTag.AddChildren(html.NewText(title))
		doc.header.AddChildren(titleTag)
	}

	body := html.NewTag("body", attrs...)
	doc.body = html.NewFragment()
	doc.Scripts = html.NewFragment()
	body.AddChildren(doc.body, doc.Scripts)

	html := html.NewTag("html")
	html.AddChildren(doc.header, body)
	doc.doc.AddChildren(html)
	return doc
}

func (d *Document) AddChildren(children ...html.Node) *Document {
	d.body.AddChildren(children...)
	return d
}

func (d *Document) Build() *Builder {
	return &Builder{
		root: d.body,
		cur:  d.body,
	}
}

func (d *Document) WriteTo(w io.Writer) (n int64, err error) {
	return d.doc.WriteTo(w)
}

func (d *Document) String() string {
	return html.String(d.doc)
}

func (d *Document) CSSLinks(hrefs ...string) *Document {
	for _, href := range hrefs {
		d.header.AddChildren(html.NewVoidTag("link", "href", href, "rel", "stylesheet", "type", "text/css"))
	}
	return d
}

func (d *Document) ScriptLinks(srcs ...string) *Document {
	for _, src := range srcs {
		d.Scripts.AddChildren(html.NewTag("script", "src", src))
	}
	return d
}

func (d *Document) Meta(attrs ...string) *Document {
	d.header.AddChildren(html.NewVoidTag("meta", attrs...))
	return d
}
