package validatebp

import (
	"bytes"
	"github.com/adamcolton/gothic/gothicgo"
	"github.com/adamcolton/gothic/gothicmodel"
	"github.com/adamcolton/gothic/gothicmodel/gomodel"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		name     string
		str      string
		expected map[string]string
	}{
		{
			name: "One Word",
			str:  "Required",
			expected: map[string]string{
				"Required": "",
			},
		},
		{
			name: "Trim one Word",
			str:  "  Required  ",
			expected: map[string]string{
				"Required": "",
			},
		},
		{
			name: "Two Words",
			str:  " Required | Foo ",
			expected: map[string]string{
				"Required": "",
				"Foo":      "",
			},
		},
		{
			name: ": in val",
			str:  "foo:bar:glorp",
			expected: map[string]string{
				"foo": "bar:glorp",
			},
		},
		{
			name: "Compound",
			str:  "Length:3-20 | Required |foo:bar:glorp",
			expected: map[string]string{
				"Length":   "3-20",
				"Required": "",
				"foo":      "bar:glorp",
			},
		},
		{
			name: "Escape",
			str:  `Invalid:;\||long\:key:fooooo`,
			expected: map[string]string{
				"Invalid":  ";|",
				"long:key": "fooooo",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := make(map[string]string)
			Parse(tc.str, v)
			assert.Equal(t, tc.expected, v)
		})
	}
}

func TestStructValidator(t *testing.T) {
	pkg := gothicgo.MustPackage("dabble")
	s := pkg.MustStruct("Imponderably")
	s.AddField("Prosperity", gothicgo.StringType)
	s.AddField("Indoline", gothicgo.IntType)

	v := Struct(s)
	v.Add("tester")
	v.AddField("Prosperity", "min-length:4|max-length:20")

	RegisterFieldWriterTo("min-length", MinLength)
	RegisterFieldWriterTo("max-length", MaxLength)

	v.ValidateMethod("Validate", true, "Prosperity")

	str := s.File().String()
	assert.Contains(t, str, "if len(i.Prosperity) < 4 && len(i.Prosperity) > 0 {")
	assert.Contains(t, str, "if len(i.Prosperity) > 20 {")
	assert.Contains(t, str, "r := validate.New()")
}

func TestModelValidator(t *testing.T) {
	m := gothicmodel.Must("Label", gothicmodel.Fields{
		{"User", "uint64"},
		{"Title", "string"},
	})
	h := ModelHelper{m}
	h.AddField("User", "min-range:100 | max-range:1000")
	h.AddField("Title", "required | min-length:3 | max-length:20")

	RegisterFieldWriterTo("min-length", MinLength)
	RegisterFieldWriterTo("max-length", MaxLength)
	RegisterFieldWriterTo("min-range", MinRange)
	RegisterFieldWriterTo("max-range", MaxRange)
	RegisterFieldWriterTo("required", Required)

	pkg := gothicgo.MustPackage("biaxillary")
	gm := gomodel.Must(pkg, m)
	buf := bytes.Buffer{}
	Model(gm).Validate(true, "Title", "User").WriteTo(&buf)
	str := buf.String()
	assert.Contains(t, str, `r.AddToField("Title", "Too Long")`)
	assert.Contains(t, str, `if l.User < 100 {`)
}
