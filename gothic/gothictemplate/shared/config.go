package shared

import (
	"github.com/adamcolton/buttress/config"
	"os"
	"path/filepath"
	"sync"
)

var setupFuncs []func()
var hasRun = false

func AddSetup(fn func()) {
	if !hasRun {
		setupFuncs = append(setupFuncs, fn)
	} else {
		fn()
	}
}

var setupOnce sync.Once

func Setup(env string) {
	setupOnce.Do(func() {
		config.Environments("dev", "test", "stage", "live")

		gopath := os.Getenv("GOPATH")
		if gopath == "" {
			gopath = os.Getenv("HOME") + "/go/"
		}
		config.SetString("gopath").As(gopath)

		importRoot := filepath.Join("github.com", "projectName")
		root := filepath.Join(gopath, "src", importRoot)

		config.SetString("path.root").As(root)
		config.SetString("path.app").As(filepath.Join(root, "app"))

		config.SetString("import.root").As(importRoot)

		// -----------Add config values here ----------------

		config.SetEnvironment(env)

		for _, setup := range setupFuncs {
			go setup()
		}
		setupFuncs = nil
		hasRun = true
	})
}
