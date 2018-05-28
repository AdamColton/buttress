package bootstrap3

import (
	"github.com/adamcolton/buttress/bootstrap3/csscontext"
	"github.com/adamcolton/buttress/html"
	"github.com/adamcolton/buttress/html/mutate"
)

// https://getbootstrap.com/docs/3.3/css/#forms

type FormStyle struct {
	HideLabels   bool
	HideFeedback bool
	Inline       bool
	Label        ColClass
	Input        ColClass
	Feedback     ColClass
}

func NewFormStyle() *FormStyle {
	return &FormStyle{
		Label:    NewColClass(),
		Input:    NewColClass(),
		Feedback: NewColClass(),
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

func (f *Form) AddInputTag(inputType, label, id string) *InputTag {
	i := NewInputTag(inputType, label, id)
	f.Inputs = append(f.Inputs, i)
	return i
}

func (f *Form) AddButton(label, icon string) *Buttons {
	b := NewButton(label, icon)
	f.Inputs = append(f.Inputs, b)
	return b
}

func (f *Form) AddButtons(label, value, icon string, context csscontext.CSSContext) *Buttons {
	b := NewButtons(label, value, icon, context)
	f.Inputs = append(f.Inputs, b)
	return b
}

func (f *Form) AddHidden(id, value string) *Hidden {
	h := NewHidden(id, value)
	f.Inputs = append(f.Inputs, h)
	return h
}

func (f *Form) AddConfirmPassword(id string) *ConfirmPassword {
	c := NewConfirmPassword(id)
	f.Inputs = append(f.Inputs, c)
	return c
}

func (f *Form) AddHTML(node html.Node) *HTML {
	h := &HTML{
		Node: node,
	}
	f.Inputs = append(f.Inputs, h)
	return h
}

func (f *Form) AddText(label, text string) *Text {
	t := &Text{
		Label: label,
		Text:  text,
	}
	f.Inputs = append(f.Inputs, t)
	return t
}

func (f *Form) AddSelect(label, id string, options []SelectOption) *Select {
	s := &Select{
		Label:   label,
		Options: options,
		ID:      id,
		Name:    id,
	}
	f.Inputs = append(f.Inputs, s)
	return s
}

func (f *Form) Post(action string) *Form {
	f.Method = "post"
	f.Action = action
	return f
}

type Input interface {
	Render() html.Node
	FormRender(formStyle *FormStyle) html.Node
	AddMutators(...mutate.Mutator)
}

type InputTag struct {
	Style *FormStyle
	Type  string
	Label string
	ID    string
	Name  string
	Value string
	mutate.MutateChain
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
	var input html.TagNode
	if i.Type == "textarea" {
		ta := html.NewTag("textarea", "class", "form-control", "id", i.ID, "name", i.Name)
		ta.AddChildren(html.NewText(i.Value))
		input = ta
	} else {
		input = html.NewVoidTag("input", "type", i.Type, "class", "form-control", "id", i.ID, "name", i.Name, "value", i.Value)
	}
	if formStyle.HideLabels {
		label.AddAttributes("class", "sr-only")
		input.AddAttributes("placeholder", i.Label)
	}
	if !formStyle.Inline {
		label.AddAttributes("class", "control-label "+formStyle.Label.String())
		div := html.NewTag("div", "class", formStyle.Input.String())
		div.AddChildren(input)
		input = div
	}
	fg := html.NewTag("div", "class", "form-group")
	fg.AddChildren(label, input)
	if !formStyle.HideFeedback {
		fg.AddChildren(html.NewTag("span", "class", "help-block "+formStyle.Feedback.String(), "id", i.ID+"-help"))
	}
	return i.MutateChain.Mutate(fg)
}

type Buttons struct {
	Style   *FormStyle
	Buttons []*Button
	mutate.MutateChain
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

func (b *Buttons) AddButton(label, value, icon string, context csscontext.CSSContext) *Button {
	btn := &Button{
		Label:   label,
		Context: context,
		Value:   value,
		Icon:    icon,
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
	return b.MutateChain.Mutate(fg)
}

type Hidden struct {
	ID    string
	Name  string
	Value string
	mutate.MutateChain
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
	hh := html.NewVoidTag("input", "type", "hidden", "id", h.ID, "name", h.Name, "value", h.Value)
	return h.MutateChain.Mutate(hh)
}

type ConfirmPassword struct {
	Style         *FormStyle
	PasswordLabel string
	AgainLabel    string
	ID            string
	Name          string
	mutate.MutateChain
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

	cph := html.NewFragment(in1.FormRender(formStyle), in2.FormRender(formStyle))
	return c.MutateChain.Mutate(cph)
}

type HTML struct {
	Node html.Node
	mutate.MutateChain
}

func (h *HTML) Render() html.Node {
	return h.Node
}

func (h *HTML) FormRender(formStyle *FormStyle) html.Node {
	return h.MutateChain.Mutate(h.Node)
}

type Text struct {
	Label string
	Text  string
	mutate.MutateChain
}

func (t *Text) Render() html.Node {
	return t.FormRender(nil)
}

func (t *Text) FormRender(formStyle *FormStyle) html.Node {
	label := html.NewTag("label")
	label.AddChildren(html.NewText(t.Label))
	var text html.Node = html.NewText(t.Text)
	if formStyle.HideLabels {
		label.AddAttributes("class", "sr-only")
		//text.AddAttributes("title", t.Label)
	}
	if !formStyle.Inline {
		label.AddAttributes("class", "control-label "+formStyle.Label.String())
		div := html.NewTag("div", "class", formStyle.Input.String())
		div.AddChildren(text)
		text = div
	}
	fg := html.NewTag("div", "class", "form-group")
	fg.AddChildren(label, text)
	return t.MutateChain.Mutate(fg)
}

type SelectOption struct {
	Value string
	Text  string
}

type Select struct {
	Label   string
	Options []SelectOption
	ID      string
	Name    string
	mutate.MutateChain
}

func NewSelect(label, id string, opts []SelectOption) *Select {
	return &Select{
		Label:   label,
		ID:      id,
		Name:    id,
		Options: opts,
	}
}

func (s *Select) Render() html.Node {
	return s.FormRender(nil)
}

func (s *Select) FormRender(formStyle *FormStyle) html.Node {
	label := html.NewTag("label")
	label.AddChildren(html.NewText(s.Label))
	var slct = html.NewTag("select", "id", s.ID, "name", s.Name, "class", "form-control")
	for _, opt := range s.Options {
		optTag := html.NewTag("option", "value", opt.Value)
		optTag.AddChildren(html.NewText(opt.Text))
		slct.AddChildren(optTag)
	}
	if formStyle.HideLabels {
		label.AddAttributes("class", "sr-only")
		slct.AddAttributes("title", s.Label)
	}
	if !formStyle.Inline {
		label.AddAttributes("class", "control-label "+formStyle.Label.String())
		div := html.NewTag("div", "class", formStyle.Input.String())
		div.AddChildren(slct)
		slct = div
	}
	fg := html.NewTag("div", "class", "form-group")
	fg.AddChildren(label, slct)
	return s.MutateChain.Mutate(fg)
}
