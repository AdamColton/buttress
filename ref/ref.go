package ref

import (
	"github.com/adamcolton/gothic/gothicgo"
	"github.com/adamcolton/gothic/gothicio"
	"github.com/adamcolton/gothic/gothicmodel/gomodel"
	"github.com/adamcolton/gothic/gothicmodel/sqlmodel"
	"strings"
	"text/template"
)

type tmplData struct {
	Name     string
	Receiver string
	Primary  string
	Call     string
}

// in FromID we're using a trick for the Zero check. The primary field on the
// reference will be initiated to the Zero val for what ever the type is, so
// we can use that to check if the primary passed in is a zero. If we're given
// a zero for the primary, return nil.
var templates = gothicio.TemplateWrapper{template.Must(template.New("templates").Parse(`
{{define "get" -}}
	{{- ""}}if {{.Receiver}}.{{.Name}} != nil {
 		return {{.Receiver}}.{{.Name}}
 	}
 	{{.Receiver}}.{{.Name}} = {{.Call}}
 	return {{.Receiver}}.{{.Name}}
{{- end}}
{{define "ToID" -}}
	if {{.Receiver}}!=nil{
		{{.Primary}} = {{.Receiver}}.{{.Primary}}
	}
	return 
{{- end}}
{{define "FromPrimary" -}}
	r := &{{.Name}}Ref{}
	if {{.Primary}}==r.{{.Primary}}{
		return nil
	}
	r.{{.Primary}} = {{.Primary}}
	return r
{{- end}}
{{define "RefMeth" -}}
	return &{{.Name}}Ref{
		{{.Primary}}:   {{.Receiver}}.{{.Primary}},
		{{.Name}}: {{.Receiver}},
	}
{{- end}}
{{define "MarshalUint" -}}
	return []byte(strconv.FormatUint({{.Receiver}}.{{.Primary}}, 10)), nil
{{- end}}
{{define "UnmarshalUint" -}}
	{{.Receiver}}.{{.Primary}}, err = strconv.ParseUint(string(b), 10, 64)
	return
{{- end}}
{{define "MarshalInt" -}}
	return []byte(strconv.FormatInt({{.Receiver}}.{{.Primary}}, 10)), nil
{{- end}}
{{define "UnmarshalInt" -}}
	{{.Receiver}}.{{.Primary}}, err = strconv.ParseInt(string(b), 10, 64)
	return
{{- end}}
{{define "MarshalString" -}}
	return []byte({{.Receiver}}.{{.Primary}}), nil
{{- end}}
{{define "UnmarshalString" -}}
	{{.Receiver}}.{{.Primary}}= string(b)
	return
{{- end}}
`))}

type EncodingType byte

const (
	NoEncoding EncodingType = iota
	UintEncoding
	IntEncoding
	StringEncoding
	beyondEncoding
)

// Builder defines how references will be generated.
type Builder struct {
	Suffix string
	SQL    string
	EncodingType
	JSON bool
	Gob  bool
}

//
func (rb Builder) New(gm *gomodel.GoModel, file *gothicgo.File, getByPrimaryFn gothicgo.FuncCaller) (*gothicgo.Struct, error) {
	var suffix string
	if rb.Suffix == "" {
		suffix = "Ref"
	} else {
		suffix = rb.Suffix
	}
	name := gm.Struct.Name() + suffix

	if file == nil {
		file = gm.Struct.File()
	}

	ref, err := file.NewStruct(name)
	if err != nil {
		return nil, err
	}

	primary := gm.GothicModel.Primary()
	primaryGoType := gomodel.Types[primary.Type()]
	primaryGoArg := gothicgo.Arg(primary.Name(), primaryGoType)

	ref.AddField(primary.Name(), primaryGoType)
	ref.Embed(gm.Struct.Ptr())

	td := tmplData{
		Primary:  primary.Name(),
		Name:     gm.Struct.Name(),
		Receiver: gm.Struct.ReceiverName,
	}

	if getByPrimaryFn != nil {
		td.Call = getByPrimaryFn.Call(file, gm.Struct.ReceiverName+"."+gm.GothicModel.Primary().Name())
		get := ref.NewMethod("Get")
		get.Body = templates.TemplateTo("get", td)
		get.Returns(gm.Struct.AsRet())
	}

	fnName := ref.Name()
	if strings.ToLower(gm.Struct.Name()) == gm.PackageRef().Name() {
		fnName = suffix
	}
	toID := file.NewFunc(fnName+"To"+primary.Name(), ref.AsArg(ref.ReceiverName))
	toID.Body = templates.TemplateTo("ToID", td)
	toID.Returns(primaryGoArg)

	// TODO: Add zero logic, if ID == 0, return nil
	fromID := file.NewFunc(fnName+"From"+primary.Name(), primaryGoArg)
	fromID.Body = templates.TemplateTo("FromPrimary", td)
	fromID.Returns(ref.AsRet())

	gomodel.Types[gm.Struct.Name()] = ref.Ptr()
	if rb.SQL != "" {
		sqlmodel.Types[gm.Struct.Name()] = rb.SQL
		sqlmodel.AddConverter(gm.Struct.Name(), toID, fromID)
	}

	refMeth := gm.NewMethod(suffix)
	refMeth.UnnamedReturns(ref.Ptr())
	refMeth.Body = templates.TemplateTo("RefMeth", td)

	if rb.EncodingType > NoEncoding && rb.EncodingType < beyondEncoding {
		if rb.JSON {
			rb.addJSON(ref, td)
		}
		if rb.Gob {
			rb.addGob(ref, td)
		}
	}

	return ref, nil
}

var (
	strconvPkg    = gothicgo.MustPackageRef("strconv")
	byteSliceType = gothicgo.ByteType.Slice()
)

func (rb Builder) addJSON(ref *gothicgo.Struct, td tmplData) {
	mj := ref.NewMethod("MarshalJSON")
	mj.UnnamedReturns(byteSliceType, gothicgo.ErrorType)
	umj := ref.NewMethod("UnmarshalJSON", gothicgo.Arg("b", byteSliceType))
	umj.Returns(gothicgo.Arg("err", gothicgo.ErrorType))

	switch rb.EncodingType {
	case UintEncoding:
		ref.File().AddRefImports(strconvPkg)
		mj.Body = templates.TemplateTo("MarshalUint", td)
		umj.Body = templates.TemplateTo("UnmarshalUint", td)
	case IntEncoding:
		ref.File().AddRefImports(strconvPkg)
		mj.Body = templates.TemplateTo("MarshalInt", td)
		umj.Body = templates.TemplateTo("UnmarshalInt", td)
	case StringEncoding:
		mj.Body = templates.TemplateTo("MarshalString", td)
		umj.Body = templates.TemplateTo("UnmarshalString", td)
	}
}

func (rb Builder) addGob(ref *gothicgo.Struct, td tmplData) {
	ge := ref.NewMethod("GobEncode")
	ge.UnnamedReturns(byteSliceType, gothicgo.ErrorType)
	gd := ref.NewMethod("GobDecode", gothicgo.Arg("b", byteSliceType))
	gd.Returns(gothicgo.Arg("err", gothicgo.ErrorType))

	switch rb.EncodingType {
	case UintEncoding:
		ref.File().AddRefImports(strconvPkg)
		ge.Body = templates.TemplateTo("MarshalUint", td)
		gd.Body = templates.TemplateTo("UnmarshalUint", td)
	case IntEncoding:
		ref.File().AddRefImports(strconvPkg)
		ge.Body = templates.TemplateTo("MarshalInt", td)
		gd.Body = templates.TemplateTo("UnmarshalInt", td)
	case StringEncoding:
		ge.Body = templates.TemplateTo("MarshalString", td)
		gd.Body = templates.TemplateTo("UnmarshalString", td)
	}
}
