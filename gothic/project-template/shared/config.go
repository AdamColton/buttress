package shared

import (
	"github.com/adamcolton/buttress/config"
	"os"
	"path/filepath"
	"sync"
)

var setupFuncs []func()

func AddSetup(fn func()) {
	setupFuncs = append(setupFuncs, fn)
}

var setupOnce sync.Once

func Setup() {
	setupOnce.Do(func() {
		gopath := os.Getenv("GOPATH")
		if gopath == "" {
			gopath = os.Getenv("HOME") + "/go/"
		}
		config.Environments("dev", "test", "prod")

		config.SetString("database.connection").
			As("devUser:devPassword@/devDatabase", "dev")

		config.SetBytes("crypto.session.auth").
			AsBase64("{{.Auth}}", "dev")

		config.SetBytes("crypto.session.enc").
			AsBase64("{{.Enc}}", "dev")

		config.SetString("port").
			As(":8080")

		for _, setup := range setupFuncs {
			setup()
		}
	})
}
