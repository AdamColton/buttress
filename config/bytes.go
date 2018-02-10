package config

import (
	"encoding/base64"
)

var decodeString = base64.URLEncoding.DecodeString
var encodeString = base64.URLEncoding.EncodeToString

var bytes = make(map[string]map[string][]byte) //[env][key] => val

type BytesSetter interface {
	As([]byte, ...string) BytesSetter
	AsBase64(val string, envs ...string) BytesSetter
}

type bytesSetter string

func (s bytesSetter) As(val []byte, envs ...string) BytesSetter {
	key := string(s)
	if len(envs) == 0 {
		envs = allEnvs
	}

	for _, envStr := range envs {
		if env, ok := bytes[envStr]; ok {
			env[key] = val
		}
	}
	return s
}

func (s bytesSetter) AsBase64(val string, envs ...string) BytesSetter {
	b, err := decodeString(val)
	if err != nil {
		panic(err)
	}
	key := string(s)
	if len(envs) == 0 {
		envs = allEnvs
	}

	for _, envStr := range envs {
		if env, ok := bytes[envStr]; ok {
			env[key] = b
		}
	}
	return s
}

func SetBytes(key string) BytesSetter { return bytesSetter(key) }
func GetBytes(key string) []byte {
	env, ok := bytes[activeEnv]
	if !ok {
		return nil
	}
	return env[key]
}
