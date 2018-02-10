package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func clear() {
	allEnvs = nil
	activeEnv = ""
	strs = make(map[string]map[string]string)
	bytes = make(map[string]map[string][]byte)
}

func TestString(t *testing.T) {
	clear()
	Environments("prod", "dev", "test")

	SetString("testString").
		As("bar", "dev").
		As("foo", "prod", "test")

	assert.Equal(t, "foo", GetString("testString"))
	SetEnvironment("dev")
	assert.Equal(t, "bar", GetString("testString"))
}

func TestBytes(t *testing.T) {
	clear()
	Environments("prod", "dev", "test")

	SetBytes("testBytes").
		As([]byte{1, 2, 3}, "dev").
		AsBase64("BAUG", "prod", "test")

	assert.Equal(t, []byte{4, 5, 6}, GetBytes("testBytes"))
	SetEnvironment("dev")
	assert.Equal(t, []byte{1, 2, 3}, GetBytes("testBytes"))
}

func TestBool(t *testing.T) {
	clear()
	Environments("prod", "dev", "test")

	SetBool("testBool").
		As(true, "dev").
		As(false, "prod", "test")

	assert.False(t, GetBool("testBool"))
	SetEnvironment("dev")
	assert.True(t, GetBool("testBool"))
}
