package main

import (
	"github.com/adamcolton/buttress/config"
	"github.com/adamcolton/gothic/gothicgo"
	"github.com/urfave/cli"
	"os"

	"projectName/blueprint/project"
	"projectName/shared"
)

func main() {
	c := cli.NewApp()
	c.Name = "Gothic Template Blueprint"
	c.Action = buildCmd

	c.Commands = []cli.Command{
		{
			Name:   "clear",
			Action: clearCmd,
		},
	}

	c.Run(os.Args)
}

func setup() {
	shared.Setup("dev")

	gothicgo.OutputPath = config.MustGetString("path.root")
	gothicgo.SetImportPath(config.MustGetString("import.root"))

	setTypes()
}

func buildCmd(c *cli.Context) error {
	setup()

	project.DeleteGeneratedFiles()
	generate()

	err := gothicgo.Export()
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func clearCmd(c *cli.Context) error {
	setup()
	project.DeleteGeneratedFiles()
	return nil
}
