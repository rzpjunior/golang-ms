package main

import (
	"fmt"
	"os"

	"git.edenfarm.id/edenlabs/cli/cmd"
)

func main() {
	if err := cmd.NewRootCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
