package bootstrap3

import (
	"github.com/adamcolton/buttress/html"
	"github.com/adamcolton/buttress/html/builder"
)

var CSSLinks = []string{
	"https://cdnjs.cloudflare.com/ajax/libs/font-awesome/4.4.0/css/font-awesome.min.css",
	"https://fonts.googleapis.com/css?family=Lato:100,300,400,700",
	"https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap.min.css",
}

var ScriptLinks = []string{
	"https://cdnjs.cloudflare.com/ajax/libs/jquery/2.1.4/jquery.min.js",
	"https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/js/bootstrap.min.js",
}

var Meta = [][]string{
	{"name", "viewport", "content", "width=device-width, initial-scale=1"},
}

type Document struct {
	*builder.Document
	nav *html.Fragment
}

func NewDocument(title string) *Document {
	nav := html.NewFragment()
	doc := builder.NewDocument(title, "id", "app-layout").
		CSSLinks(CSSLinks...).
		ScriptLinks(ScriptLinks...).
		AddChildren(nav)

	for _, m := range Meta {
		doc.Meta(m...)
	}

	return &Document{
		Document: doc,
		nav:      nav,
	}
}

func (d *Document) CSSLinks(hrefs ...string) *Document {
	d.Document.CSSLinks(hrefs...)
	return d
}

func (d *Document) ScriptLinks(srcs ...string) *Document {
	d.Document.ScriptLinks(srcs...)
	return d
}
