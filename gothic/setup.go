package main

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/urfave/cli"
	"os"
	"text/template"
)

var dirs = []string{
	"app",
	"app/models",
	"app/tools",
	"app/resources",
	"blueprint",
	"blueprint/models",
	"shared",
}

func setup(c *cli.Context) error {
	for _, dir := range dirs {
		os.Mkdir(dir, 0766)
	}
	generateConfig()
	generateDB()
	return nil
}

var configTmpl = template.Must(template.New("config").Parse(`package shared

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
`))

func generateConfig() {
	file, err := os.Create("shared/config.go")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	auth := make([]byte, 32)
	enc := make([]byte, 32)
	rand.Read(auth)
	rand.Read(enc)
	var data = struct {
		Auth, Enc string
	}{
		Auth: base64.URLEncoding.EncodeToString(auth),
		Enc:  base64.URLEncoding.EncodeToString(enc),
	}
	configTmpl.Execute(file, data)
}

var dbTmpl = `package shared

import (
	"github.com/adamcolton/buttress/config"
	"github.com/adamcolton/buttress/gsql"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	AddSetup(setupDB)
}

func setupDB() {
	_, err := gsql.SetConn("mysql", config.GetString("database.connection"))
	if err != nil {
		panic(err)
	}
}
`

func generateDB() {
	file, err := os.Create("shared/db.go")
	if err != nil {
		panic(err)
	}
	file.WriteString(dbTmpl)
	file.Close()
}
