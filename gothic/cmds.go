package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/adamcolton/buttress/gothic/gothicproject"
	"github.com/urfave/cli"
	"os/exec"
	"path/filepath"
	"time"
)

func loadGPAndOuputGen(force bool) (*gothicproject.GothicProject, error) {
	gp, err := gothicproject.Load()
	if err != nil {
		return nil, err
	}

	out, err := gp.RunBlueprint(force)
	if out != "" {
		fmt.Print(out)
	}

	return gp, err
}

func blueprintCmd(c *cli.Context) error {
	_, err := loadGPAndOuputGen(true)
	return err
}

func testsCmd(c *cli.Context) error {
	gp, err := loadGPAndOuputGen(false)
	if err != nil {
		return err
	}

	for _, td := range gp.Tests {
		ftd := fmt.Sprintf("./%s", filepath.Join(gp.App, td))
		out, err := exec.Command("go", "test", ftd).CombinedOutput()
		if err != nil {
			return err
		}
		fmt.Print(string(out))
	}

	return nil
}

func runCmd(c *cli.Context) error {
	gp, err := loadGPAndOuputGen(false)
	if err != nil {
		return err
	}

	files, err := filepath.Glob(filepath.Join(gp.App, "*.go"))
	if err != nil {
		return err
	}

	args := append([]string{"run"}, files...)
	args = append(args, c.Args()...)
	runInteractive("go", args...)
	return nil
}

func keyCmd(c *cli.Context) error {
	k := make([]byte, 32)
	rand.Read(k)
	fmt.Println(base64.URLEncoding.EncodeToString(k))
	return nil
}

func getTimeCmd(c *cli.Context) error {
	f := gothicproject.TimeFormat
	if gp, _ := gothicproject.Load(); gp != nil {
		f = gp.TimeFormat
	}
	fmt.Println(time.Now().Format(f))
	return nil
}
