package validatebp

import (
	"github.com/adamcolton/gothic/gothicmodel"
	"github.com/adamcolton/gothic/gothicmodel/gomodel"
)

var ValidateKey = "validate"

// Model takes a GoModel and returns a StructValidator that references the
// GoModel struct using the Meta data on the GothicModel.
func Model(goModel *gomodel.GoModel) *StructValidator {
	sv := &StructValidator{
		Struct:          goModel.Struct,
		validators:      make(map[string]string),
		fieldValidators: make(map[string]map[string]string),
	}

	if v, ok := goModel.GothicModel.Meta(ValidateKey); ok {
		sv.Add(v)
	}

	for _, f := range goModel.GothicModel.Fields() {
		if v, ok := f.Meta(ValidateKey); ok {
			sv.AddField(f.Name(), v)
		}
	}

	return sv
}

type ModelHelper struct {
	GothicModel *gothicmodel.GothicModel
}

func (m ModelHelper) Add(rule string) {
	m.GothicModel.AddMeta(ValidateKey, rule)
}

func (m ModelHelper) AddField(field, rule string) bool {
	f, ok := m.GothicModel.Field(field)
	if !ok {
		return false
	}
	f.AddMeta(ValidateKey, rule)
	return true
}
