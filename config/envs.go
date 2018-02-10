package config

var allEnvs []string

func Environments(environments ...string) {
	for _, env := range environments {
		strs[env] = make(map[string]string)
		bytes[env] = make(map[string][]byte)
		bools[env] = make(map[string]bool)
	}
	if activeEnv == "" && len(environments) > 0 {
		activeEnv = environments[0]
	}
	allEnvs = append(allEnvs, environments...)
}

var activeEnv string

func SetEnvironment(env string) {
	activeEnv = env
}
