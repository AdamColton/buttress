package bootstrap3

import (
	"github.com/adamcolton/buttress/bootstrap3/csscontext"
	"github.com/adamcolton/buttress/html"
)

// https://getbootstrap.com/docs/3.3/css/#forms

type FormStyle struct {
	HideLabels bool
	Inline     bool
	Label      ColClass
	Input      ColClass
}

func NewFormStyle() *FormStyle {
	return &FormStyle{
		Label: NewColClass(),
		Input: NewColClass(),
	}
}

type Form struct {
	Style  *FormStyle
	Action string
	Method string
	Inputs []Input
}

func NewForm(formStyle *FormStyle) *Form {
	if formStyle == nil {
		formStyle = NewFormStyle()
	}
	return &Form{
		Style: formStyle,
	}
}

func (f *Form) Render() html.Node {
	class := "form-horizontal"
	if f.Style.Inline {
		class = "form-inline"
	}
	tag := html.NewTag("form", "method", f.Method, "action", f.Action, "class", class)
	for _, i := range f.Inputs {
		tag.AddChildren(i.FormRender(f.Style))
	}
	return tag
}

func (f *Form) AddInputs(inputs ...Input) {
	f.Inputs = append(f.Inputs, inputs...)
}

func (f *Form) InputTag(inputType, label, id string) *InputTag {
	i := NewInputTag(inputType, label, id)
	f.Inputs = append(f.Inputs, i)
	return i
}

func (f *Form) Button(label, icon string) *Buttons {
	b := NewButton(label, icon)
	f.Inputs = append(f.Inputs, b)
	return b
}

func (f *Form) Buttons(label, value, icon string, context csscontext.CSSContext) *Buttons {
	b := NewButtons(label, value, icon, context)
	f.Inputs = append(f.Inputs, b)
	return b
}

func (f *Form) Hidden(id, value string) *Hidden {
	h := NewHidden(id, value)
	f.Inputs = append(f.Inputs, h)
	return h
}

func (f *Form) ConfirmPassword(id string) *ConfirmPassword {
	c := NewConfirmPassword(id)
	f.Inputs = append(f.Inputs, c)
	return c
}

func (f *Form) Post(action string) {
	f.Method = "post"
	f.Action = action
}

type Input interface {
	Render() html.Node
	FormRender(formStyle *FormStyle) html.Node
}

type InputTag struct {
	Style *FormStyle
	Type  string
	Label string
	ID    string
	Name  string
	Value string
}

func NewInputTag(inputType, label, id string) *InputTag {
	return &InputTag{
		Type:  inputType,
		Label: label,
		ID:    id,
		Name:  id,
	}
}

func (i *InputTag) Render() html.Node {
	return i.FormRender(nil)
}

func (i *InputTag) FormRender(formStyle *FormStyle) html.Node {
	if i.Style != nil {
		formStyle = i.Style
	}
	label := html.NewTag("label", "for", i.ID)
	label.AddChildren(html.NewText(i.Label))
	var input html.TagNode = html.NewVoidTag("input", "type", i.Type, "class", "form-control", "id", i.ID, "name", i.Name, "value", i.Value)
	if !formStyle.Inline {
		label.AddAttributes("class", "control-label "+formStyle.Label.String())
		div := html.NewTag("div", "class", formStyle.Input.String())
		div.AddChildren(input)
		input = div
	}
	fg := html.NewTag("div", "class", "form-group")
	fg.AddChildren(label, input)
	return fg
}

type Buttons struct {
	Style   *FormStyle
	Buttons []*Button
}

type Button struct {
	Label   string
	Context csscontext.CSSContext
	Value   string
	Icon    string
}

func NewButton(label, icon string) *Buttons {
	return &Buttons{
		Buttons: []*Button{
			&Button{
				Label:   label,
				Context: csscontext.Primary(),
				Icon:    icon,
			},
		},
	}
}

func NewButtons(label, value, icon string, context csscontext.CSSContext) *Buttons {
	return &Buttons{
		Buttons: []*Button{
			&Button{
				Label:   label,
				Context: context,
				Value:   value,
				Icon:    icon,
			},
		},
	}
}

func (b *Buttons) AddButton(label, value string, context csscontext.CSSContext) *Button {
	btn := &Button{
		Label:   label,
		Context: context,
		Value:   value,
	}
	b.Buttons = append(b.Buttons, btn)
	return btn
}

func (b *Buttons) Render() html.Node {
	return b.FormRender(nil)
}

func (b *Buttons) FormRender(formStyle *FormStyle) html.Node {
	if b.Style != nil {
		formStyle = b.Style
	}
	fg := html.NewTag("div", "class", "form-group")
	buttons := make([]html.Node, len(b.Buttons))
	for i, btn := range b.Buttons {
		bt := html.NewTag("button", "class", "btn btn-"+btn.Context.String(), "value", btn.Value)
		Icon(bt, btn.Icon)
		bt.AddChildren(html.NewText(btn.Label))
		buttons[i] = bt
	}
	if !formStyle.Inline {
		fg.AddChildren(html.NewTag("div", "class", formStyle.Label.String()))
		div := html.NewTag("div", "class", formStyle.Input.String())
		div.AddChildren(buttons...)
		buttons = []html.Node{div}
	}
	fg.AddChildren(buttons...)
	return fg
}

type Hidden struct {
	ID    string
	Name  string
	Value string
}

func NewHidden(id, value string) *Hidden {
	return &Hidden{
		ID:    id,
		Name:  id,
		Value: value,
	}
}

func (h *Hidden) Render() html.Node {
	return h.FormRender(nil)
}

func (h *Hidden) FormRender(formStyle *FormStyle) html.Node {
	return html.NewVoidTag("input", "type", "hidden", "id", h.ID, "name", h.Name, "value", h.Value)
}

type ConfirmPassword struct {
	Style         *FormStyle
	PasswordLabel string
	AgainLabel    string
	ID            string
	Name          string
}

func NewConfirmPassword(id string) *ConfirmPassword {
	return &ConfirmPassword{
		PasswordLabel: "Password",
		AgainLabel:    "Again",
		ID:            id,
		Name:          id,
	}
}

func (c *ConfirmPassword) Render() html.Node {
	return c.FormRender(nil)
}

func (c *ConfirmPassword) FormRender(formStyle *FormStyle) html.Node {
	if c.Style != nil {
		formStyle = c.Style
	}
	in1 := NewInputTag("password", c.PasswordLabel, c.ID)
	in1.Name = c.Name
	in2 := NewInputTag("password", c.AgainLabel, c.ID+"_confirm")
	in2.Name = c.Name + "_confirm"

	return html.NewFragment(in1.FormRender(formStyle), in2.FormRender(formStyle))
}
