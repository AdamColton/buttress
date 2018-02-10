package ref

import (
	"bytes"
	"github.com/adamcolton/gothic/gothicgo"
	"github.com/adamcolton/gothic/gothicmodel"
	"github.com/adamcolton/gothic/gothicmodel/gomodel"
	"github.com/adamcolton/gothic/gothicmodel/sqlmodel"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

type wc struct {
	*bytes.Buffer
}

type testFuncCaller struct{}

func (testFuncCaller) Call(pre gothicgo.Prefixer, args ...string) string {
	t := gothicgo.MustPackageRef("tester")
	return pre.Prefix(t) + "GetFooByID(" + strings.Join(args, ", ") + ")"
}

func (testFuncCaller) Args() []gothicgo.NameType {
	return nil
}

func (testFuncCaller) Rets() []gothicgo.NameType {
	return nil
}

func (*wc) Close() error { return nil }

func TestBuilder(t *testing.T) {
	rb := Builder{
		SQL:          "int unsigned",
		EncodingType: StringEncoding,
		JSON:         true,
		Gob:          true,
	}
	user := gothicmodel.Must("User", gothicmodel.Fields{
		{"Key", "uint64"},
		{"LastLogin", "datetime"},
		{"Name", "string"},
		{"Password", "password"},
	})
	pkg := gothicgo.MustPackage("user")
	gm, err := gomodel.New(pkg, user)
	assert.NoError(t, err)

	ref, err := rb.New(gm, nil, testFuncCaller{})
	assert.NoError(t, err)
	assert.NotNil(t, ref)
	buf := &bytes.Buffer{}
	gm.Struct.File().Writer = &wc{buf}
	gm.Struct.File().Prepare()
	gm.Struct.File().Generate()
	s := buf.String()
	assert.Contains(t, s, "func RefToKey(u *UserRef) (Key uint64) {")
	assert.Contains(t, s, "func RefFromKey(Key uint64) *UserRef {")
	assert.Contains(t, s, "func (u *User) Ref() *UserRef {")
	assert.Contains(t, s, "u.User, _ = tester.GetFooByID(u.Key)")
	assert.Contains(t, s, "u.User")

	s = ref.String()
	assert.Regexp(t, "Key +uint64", s)
	assert.Contains(t, s, "type UserRef struct {")

	// confirm that "User" type has been added
	data := gothicmodel.Must("Data", gothicmodel.Fields{
		{"Key", "uint64"},
		{"User", "User"},
		{"Data", "string"},
	})
	gm, err = gomodel.New(pkg, data)
	assert.NoError(t, err)
	s = gm.String()
	assert.Regexp(t, `User +\*UserRef`, s)
	sm := sqlmodel.New(gm)
	s = sm.Scanner().String()
	assert.Contains(t, s, "var User uint64")
	assert.Contains(t, s, "d.User = RefFromKey(User)")
	assert.Contains(t, s, "err := row.Scan(&(d.Key), &(User), &(d.Data))")
}
