package config

import (
	"fmt"
	"sort"
	"strings"
)

const (
	_ = iota
	strKind
	boolKind
	bytesKind
	intKind
	floatKind
)

type key struct {
	env, name string
	kind      int
}

var (
	allenvironments []string
	activeEnv       string
	vals            = make(map[key]interface{})
)

func Environments(environments ...string) {
	if activeEnv == "" && len(environments) > 0 {
		activeEnv = environments[0]
	}
	allenvironments = append(allenvironments, environments...)
}

func SetEnvironment(env string) {
	activeEnv = env
}

func String(environments ...string) string {
	if len(environments) == 0 {
		environments = []string{activeEnv}
	}
	var out []string
	for _, env := range environments {
		out = append(out, fmt.Sprint("== ", env, " =="))
		var valStrs []string
		for k, v := range vals {
			if k.env == env {
				valStrs = append(valStrs, fmt.Sprintf("  %s: %s", k.name, format(v)))
			}
		}
		sort.Slice(valStrs, func(i, j int) bool { return valStrs[i] < valStrs[j] })
		out = append(out, strings.Join(valStrs, "\n"))
	}

	return strings.Join(out, "\n")
}

func format(v interface{}) string {
	switch t := v.(type) {
	case string:
		return fmt.Sprintf("%q", t)
	case []byte:
		out := make([]string, len(t)+2)
		out[0] = "["
		out[len(t)+1] = "]"
		for i, b := range t {
			out[i+1] = fmt.Sprintf("%02x", b)
		}
		return strings.Join(out, " ")
	}
	return fmt.Sprintf("%#v", v)
}
