package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func clear() {
	allenvironments = nil
	activeEnv = ""
	vals = make(map[key]interface{})
}

func TestString(t *testing.T) {
	clear()
	Environments("prod", "dev", "test")

	SetString("testString").
		As("bar", "dev").
		As("foo", "prod", "test")

	assert.Equal(t, "foo", MustGetString("testString"))
	SetEnvironment("dev")
	assert.Equal(t, "bar", MustGetString("testString"))

	_, err := GetString("notSet")
	if assert.Error(t, err) {
		assert.Equal(t, "Cound not find string 'notSet' in environment 'dev'", err.Error())
	}
}

func TestBytes(t *testing.T) {
	clear()
	Environments("prod", "dev", "test")

	SetBytes("testBytes").
		As([]byte{1, 2, 3}, "dev").
		AsBase64("BAUG", "prod", "test")

	assert.Equal(t, []byte{4, 5, 6}, MustGetBytes("testBytes"))
	SetEnvironment("dev")
	assert.Equal(t, []byte{1, 2, 3}, MustGetBytes("testBytes"))
}

func TestBool(t *testing.T) {
	clear()
	Environments("prod", "dev", "test")

	SetBool("testBool").
		As(true, "dev").
		As(false, "prod", "test")

	assert.False(t, MustGetBool("testBool"))
	SetEnvironment("dev")
	assert.True(t, MustGetBool("testBool"))
}

func TestInt(t *testing.T) {
	clear()
	Environments("prod", "dev", "test")

	SetInt("testInt").
		As(22, "dev").
		As(55, "prod", "test")

	assert.Equal(t, 55, MustGetInt("testInt"))
	SetEnvironment("dev")
	assert.Equal(t, 22, MustGetInt("testInt"))
}

func TestFloat(t *testing.T) {
	clear()
	Environments("prod", "dev", "test")

	SetFloat("testFloat").
		As(3.14159, "dev").
		As(2.71828, "prod", "test")

	assert.Equal(t, 2.71828, MustGetFloat("testFloat"))
	SetEnvironment("dev")
	assert.Equal(t, 3.14159, MustGetFloat("testFloat"))
}

func TestOuputString(t *testing.T) {
	clear()
	Environments("prod", "dev", "test")
	SetString("tosString").
		As("bar", "dev").
		As("foo", "prod", "test")
	SetBytes("tosBytes").
		As([]byte{9, 15, 200}, "dev").
		AsBase64("BAUG", "prod", "test")
	SetBool("tosBool").
		As(true, "dev").
		As(false, "prod", "test")

	expected := `== prod ==
  tosBool: false
  tosBytes: [ 04 05 06 ]
  tosString: "foo"
== dev ==
  tosBool: true
  tosBytes: [ 09 0f c8 ]
  tosString: "bar"
== test ==
  tosBool: false
  tosBytes: [ 04 05 06 ]
  tosString: "foo"`

	assert.Equal(t, expected, String("prod", "dev", "test"))

	expected = `== prod ==
  tosBool: false
  tosBytes: [ 04 05 06 ]
  tosString: "foo"`
	assert.Equal(t, expected, String())
}
