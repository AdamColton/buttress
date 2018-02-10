package main

import (
	"github.com/urfave/cli"
	"os"
)

func main() {
	c := cli.NewApp()
	c.Name = "Gothic"
	c.Usage = "Tool for managing a Gothic project"

	c.Commands = []cli.Command{
		{
			Name:   "blueprint",
			Action: blueprintCmd,
			Usage:  "Forces the blueprint to run",
		},
		{
			Name:   "test",
			Action: testsCmd,
			Usage:  "Runs the tests listed in gothic.json. Will run blueprint if it has changed.",
		},
		{
			Name:   "run",
			Action: runCmd,
			Usage:  "Runs the app. Will run blueprint if it has changed.",
		},
		{
			Name:  "gen",
			Usage: "Generates useful data",
			Subcommands: []cli.Command{
				{
					Name:   "key",
					Action: keyCmd,
					Usage:  "Reads 32 bytes from crypto/rand and writes them as a URLEncoding base64 string",
				},
				{
					Name:   "time",
					Action: getTimeCmd,
					Usage:  "Outputs the date and time down to the minute, useful for migration names",
				},
			},
		},
	}

	c.Run(os.Args)
}
