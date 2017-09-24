package bootstrap3

import (
	"github.com/adamcolton/buttress/html"
	"github.com/adamcolton/buttress/html/builder"
	"github.com/adamcolton/buttress/html/qhtml"
	"github.com/adamcolton/buttress/html/query"
)

var standardNav = query.NewPrototype(qhtml.MustParse(`
	<nav.navbar.navbar-default.navbar-static-top>
		<div.container>
			<div.navbar-header>
				<button.navbar-toggle.collapsed type="button" data-toggle="collapse" data-target="#navbar" aria-expanded="false" aria-controls="navbar">
					<span.sr-only>Toggle navigation</>
					<span.icon-bar></>
					<span.icon-bar></>
					<span.icon-bar></>
				</>
				<a#brand.navbar-brand></>
			</>
			<div#navbar.collapse.navbar-collapse>
				<ul#nav-left.nav.navbar-nav></>
				<ul#nav-right.nav.navbar-nav.navbar-right></>
			</>
		</>
	</>
`)).Queries(
	"brand", "#brand",
	"left", "#nav-left",
	"right", "#nav-right",
)

func (d *Document) StandardNav(nav *Nav) {
	newNav := standardNav.Clone()
	brand := standardNav.Query("brand", newNav).(*html.Tag)
	if nav.Href != "" {
		brand.AddAttributes("href", nav.Href)
	}
	if nav.Name != "" {
		brand.AddChildren(html.NewText(nav.Name))
	}

	standardMenu{nav.Menu(Left)}.populate(standardNav.Query("left", newNav).(*html.Tag))
	standardMenu{nav.Menu(Right)}.populate(standardNav.Query("right", newNav).(*html.Tag))

	d.nav.RemoveAll()
	d.nav.AddChildren(newNav)
}

type standardMenu struct {
	Menu
}

func (menu standardMenu) populate(container html.ContainerNode) {
	for _, mi := range menu.Menu {
		b := builder.New()
		if mi.HasSub() {
			b.Tag("li", "class", "dropdown").
				Tag("a", mi.Attrs("class", "dropdown-toggle", "data-toggle", "dropdown", "role", "button", "aria-haspopup", "true", "aria-expanded", "false")...)
			Icon(b, mi.Icon())
			b.Text(mi.Title()).
				Tag("span", "class", "caret").
				Close().Close().
				Tag("ul", "class", "dropdown-menu")
			for _, smi := range mi.SubMenu() {
				standardMenuItem{smi}.addSubmenu(b)
			}
			b.Close()
		} else {
			b.Tag("li").
				Tag("a", mi.Attrs("role", "button")...)
			Icon(b, mi.Icon())
			b.Text(mi.Title()).
				Close()
		}
		container.AddChildren(b.Cur())
	}
}

type standardMenuItem struct {
	MenuItem
}

func (smi standardMenuItem) addSubmenu(b *builder.Builder) {
	b.Tag("li").
		Tag("a", smi.Attrs()...)
	Icon(b, smi.Icon())
	b.Text(smi.Title()).
		Close().Close()
}
