package config

var bools = make(map[string]map[string]bool) //[env][key] => val

type BoolSetter interface {
	As(bool, ...string) BoolSetter
}

type boolSetter string

func (s boolSetter) As(val bool, envs ...string) BoolSetter {
	key := string(s)
	if len(envs) == 0 {
		envs = allEnvs
	}

	for _, envStr := range envs {
		if env, ok := bools[envStr]; ok {
			env[key] = val
		}
	}
	return s
}

func SetBool(key string) BoolSetter { return boolSetter(key) }
func GetBool(key string) bool {
	env, ok := bools[activeEnv]
	if !ok {
		println("not found")
		return false
	}
	return env[key]
}
