package main

import (
	"github.com/adamcolton/buttress/bootstrap3"
	"github.com/adamcolton/buttress/bootstrap3/bootstrap3model"
	"github.com/adamcolton/buttress/bootstrap3/csscontext"
	"github.com/adamcolton/buttress/bootstrap3/csssize"
	"github.com/adamcolton/buttress/html"
	"github.com/adamcolton/gothic/gothicmodel"
	"github.com/adamcolton/socketServer"
	"net/http"
)

var (
	panelClass      bootstrap3.ColClass
	formStyle       *bootstrap3.FormStyle
	formBuilder     *bootstrap3model.FormBuilder
	userModel       *gothicmodel.Model
	confirmPassword = bootstrap3model.TypeMap{"password": bootstrap3model.ConfirmPassword}
)

func main() {
	panelClass = bootstrap3.NewColClass()
	panelClass.SetSize(12, csssize.ExtraSmall())
	panelClass.SetSizeOffset(10, 1, csssize.Medium())

	formStyle = bootstrap3.NewFormStyle()
	formStyle.Label.SetSizeOffset(3, 2, csssize.ExtraSmall())
	formStyle.Input.SetSize(4, csssize.ExtraSmall())

	formBuilder = bootstrap3model.NewFormBuilder(formStyle)

	userModel = gothicmodel.Must("User", gothicmodel.Fields{
		{"ID", "uint"},
		{"Name", "string"},
		{"Password", "password"},
		{"Age", "int"},
		{"LastLogin", "datetime"},
	})

	s := socketServer.New()
	s.HandleFunc("/", home)
	s.HandleFunc("/login", login)
	s.HandleFunc("/create", create)
	s.HandleFunc("/panels", panels)
	s.ListenAndServe(":6060")
}

func page() *bootstrap3.Document {
	nav := bootstrap3.NewNav("Gothic Testing", "/")
	demos := nav.Add(bootstrap3.Right, "demos", "Demos", "")
	demos.Sub("create", "New User", "user-plus", "/create")
	demos.Sub("login", "Login", "sign-in", "/login")
	demos.Sub("panels", "Panels", "", "/panels")
	nav.Add(bootstrap3.Left, "home", "Home", "home").SetHref("/")

	d := bootstrap3.NewDocument("Bootstrap Demo")
	d.StandardNav(nav)
	return d
}

func panel(title string, context csscontext.CSSContext, contents html.Node) *bootstrap3.Container {
	panel := bootstrap3.NewPanel(title, contents)
	panel.ContextClass = context

	c := bootstrap3.NewContainer()
	c.AddNodes(panelClass, panel.Render())
	return c
}

func home(w http.ResponseWriter, r *http.Request) {
	d := page()
	p := panel("Home", csscontext.Warning(), html.NewText("Welcome to the Gothic Bootstrap 3 demo."))
	d.AddChildren(p.Render())
	d.Write(w)
}

func create(w http.ResponseWriter, r *http.Request) {
	f := formBuilder.Form(userModel.SkipFields("ID"), confirmPassword)
	f.Button("Create")

	d := page()
	p := panel("New User", csscontext.Primary(), f.Render())
	d.AddChildren(p.Render())
	d.Write(w)
}

func login(w http.ResponseWriter, r *http.Request) {
	f := formBuilder.Form(userModel.Fields("Name", "Password"), nil)
	f.Button("Login")

	d := page()
	p := panel("Login", csscontext.Primary(), f.Render())
	d.AddChildren(p.Render())
	d.Write(w)
}

func panels(w http.ResponseWriter, r *http.Request) {
	p1 := bootstrap3.NewPanel("Yesterday", html.NewText("stuff I did yesterday"))
	p2 := bootstrap3.NewPanel("Today", html.NewText("stuff I'm doing today"))
	p3 := bootstrap3.NewPanel("Tomorrow", html.NewText("stuff I'll do tomorrow"))

	cc := bootstrap3.NewColClass()
	cc.SetSize(4, csssize.Small())

	c := bootstrap3.NewContainer()
	c.AddNodes(cc, p1.Render(), p2.Render(), p3.Render())
	d := page()
	d.AddChildren(c.Render())
	d.Write(w)
}
