package builder

import (
	"github.com/adamcolton/buttress/html"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuilder(t *testing.T) {
	div := New().
		Tag("div", "id", "testing").
		Text("this is a test").
		Root()
	expected := `<div id="testing">this is a test</div>`
	assert.Equal(t, expected, html.String(div))
}

func TestDocument(t *testing.T) {
	html.Padding = "  "
	doc := NewDocument("Test Doc", "id", "body").
		CSSLinks("/css/foo.css", "/css/bar.css").
		ScriptLinks("/js/foo.js", "/js/bar.js")

	doc.Build().
		Tag("p").Text("I am a paragraph").Close().
		Tag("ul").
		Tag("li").Text("Item 1").Close().
		Tag("li").Text("Item 2").Close().
		Tag("li").Text("Item 3").Close()

	expected := `<!DOCTYPE html>
<html>
  <header>
    <title>Test Doc</title>
    <link href="/css/foo.css" rel="stylesheet" type="text/css" />
    <link href="/css/bar.css" rel="stylesheet" type="text/css" />
  </header>
  <body id="body">
    <p>I am a paragraph</p>
    <ul>
      <li>Item 1</li>
      <li>Item 2</li>
      <li>Item 3</li>
    </ul>
    <script src="/js/foo.js"></script>
    <script src="/js/bar.js"></script>
  </body>
</html>`

	assert.Equal(t, expected, doc.String())
}
