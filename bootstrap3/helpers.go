package bootstrap3

import (
	"github.com/adamcolton/buttress/html"
)

type ChildrenAdder interface {
	AddChildren(...html.Node)
}

func Icon(ca ChildrenAdder, icon string) {
	if icon != "" {
		ca.AddChildren(html.NewTag("i", "class", "fa fa-"+icon, "aria-hidden", "true"))
	}
}
