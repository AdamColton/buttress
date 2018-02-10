package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
)

/*
This will eventually become part of the gothic project and will be a general
tool for managing a gothic project. It can run the blueprint or the app or build
the app for deploy.
*/

func main() {
	c := cli.NewApp()
	c.Name = "Gothic"
	c.Action = runApp

	c.Commands = []cli.Command{
		{
			Name:   "gen",
			Action: generate,
		},
		{
			Name:   "setup",
			Action: setup,
		},
	}

	c.Run(os.Args)
}

func runApp(c *cli.Context) error {
	fmt.Println(os.Getwd())
	return nil
}

func generate(c *cli.Context) error {
	fmt.Println("nothing")
	return nil
}
