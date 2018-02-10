package main

import (
	"fmt"
	"os"
	"os/exec"
)

func runInteractive(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Start(); err != nil {
		fmt.Println(err)
	}
	cmd.Wait()
}
