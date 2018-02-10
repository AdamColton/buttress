package config

var strs = make(map[string]map[string]string) //[env][key] => val

type StrSetter interface {
	As(string, ...string) StrSetter
}

type strSetter string

func (s strSetter) As(val string, envs ...string) StrSetter {
	key := string(s)
	if len(envs) == 0 {
		envs = allEnvs
	}

	for _, envStr := range envs {
		if env, ok := strs[envStr]; ok {
			env[key] = val
		}
	}
	return s
}

func SetString(key string) StrSetter { return strSetter(key) }
func GetString(key string) string {
	env, ok := strs[activeEnv]
	if !ok {
		return ""
	}
	return env[key]
}
