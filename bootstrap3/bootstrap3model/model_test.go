package bootstrap3model

import (
	"github.com/adamcolton/buttress/bootstrap3"
	"github.com/adamcolton/buttress/bootstrap3/csssize"
	"github.com/adamcolton/buttress/html"
	"github.com/adamcolton/gothic/gothicmodel"
	"github.com/stretchr/testify/assert"
	"testing"
)

func LapizValueBuilder(name string) string { return "$" + name }

func setup() (*gothicmodel.GothicModel, *bootstrap3.FormStyle) {
	m, err := gothicmodel.New("test", gothicmodel.Fields{
		{"ID", "uint"},
		{"Name", "string"},
		{"Password", "password"},
		{"Age", "int"},
		{"LastLogin", "datetime"},
	})
	if err != nil {
		panic(err)
	}

	fs := bootstrap3.NewFormStyle()
	fs.Label.SetSizeOffset(3, 2, csssize.ExtraSmall())
	fs.Input.SetSize(4, csssize.ExtraSmall())

	return m, fs
}

func TestModel(t *testing.T) {
	// simulate a login
	m, fs := setup()
	assert.NotNil(t, m)
	b := NewFormBuilder(fs)
	b.ValueBuilder = LapizValueBuilder

	f := b.NewForm(m.Fields("Name", "Password"), nil)
	f.AddButton("Save", "")

	// simulate a create
	f = b.NewForm(m.Fields("Name", "Password"), TypeMap{"password": ConfirmPassword})
	f.AddButton("Create", "")
	s := html.String(f.Render())
	println(s)
}
